package goutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/AndrewDonelson/golog"
)

const (
	// Prevent infinite loops
	MAX_ITERATIONS = 100
)

// Parsed{} contains each JSON element, remembering where it was found
// - DistinctName is shortest unique name across all filenames
// - Position is element # within the file: 0 if single element, 1...N if array of N elements
type Parsed struct {
	FileName     string
	DistinctName string
	Position     int
	ElementMap   ElementMap
}

// ElementMap is the JSON element parsed into a key-value map
type ElementMap map[string]interface{}

// When mutiple files are parsed, a field in each element is specified as the Id
// - This element Id is used as the ParsedMap key (so becomes a required field)
type ParsedMap map[string][]Parsed

// The Parsed map form is then collapsed into a single data object result per Id
type ResultMap map[string]interface{}

// Read a single config file, return a struct, where 'data' is a pointer to that struct
func ReadConfigFile(data interface{}, filename string) (err error) {
	var b []byte
	var errList ErrList
	var config interface{}
	var k string

	// Make sure data is a pointer to a struct
	k = reflect.TypeOf(data).Kind().String()
	if k != "ptr" {
		err = fmt.Errorf("ReadConfigFile: 'data' must be ptr, not %s", k)
		panic(err)
	}
	st := reflect.TypeOf(data).Elem()
	sv := reflect.ValueOf(data).Elem()

	_, _, err = ValidateFile(filename)
	if err != nil {
		return
	}

	if b, err = ioutil.ReadFile(filename); err != nil {
		err = fmt.Errorf("reading file: %v", err)
		return
	}

	// Parse config file into map[string]interface{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		err = fmt.Errorf("%v [%s]", err, filename)
		return
	}

	// Each file can contain a single element of type 'data'
	k = reflect.TypeOf(config).Kind().String()
	if k != "map" {
		err = fmt.Errorf("contains type %q, must be a JSON element [%s]", k, filename)
		return
	}

	golog.Log.Debugf("Parsing single element [%s]", filename)
	parsed := Parsed{
		FileName:     filename,
		DistinctName: filepath.Base(filename),
		ElementMap:   config.(map[string]interface{}),
	}

	// store in resultMap
	parsedMap := make(ParsedMap)
	parsedMap["default"] = []Parsed{parsed}

	// Parse dataMap entries into data object (st, sv) fields
	parseConfig(st, sv, parsedMap, &errList)

	// Compile error list into an error message
	err = errList.Get()
	return
}

// Read a list of config files into a map of structs, where 'data' points to struct and idName is field for map key
// - Can configure an application using one or more JSON files
// - For example, put general settings in one file, credentials in a second file.
func ReadConfigFiles(data interface{}, idName string, filenames ...string) (resultMap ResultMap, err error) {
	var b []byte
	var errList ErrList
	var config, v interface{}
	var file FileDetail
	var fileDetails []FileDetail
	var parsedMap ParsedMap
	var parsedArr []Parsed
	var k, elementId, idTag string
	var i int
	var ok, tag bool

	// Make sure data is a pointer to a struct
	k = reflect.TypeOf(data).Kind().String()
	if k != "ptr" {
		err = fmt.Errorf("ReadConfigFiles: 'data' must be ptr, not %s", k)
		panic(err)
	}
	st := reflect.TypeOf(data).Elem()
	sv := reflect.ValueOf(data).Elem()

	// Locate field idName and `json:"idTag"`
	for i = 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if field.Name == idName {
			idTag, tag = field.Tag.Lookup("json")
			ok = true
			break
		}
	}
	if !ok {
		err = fmt.Errorf("ReadConfigFiles: 'data' does not contain field %s", idName)
		panic(err)
	}

	// Validate file list
	fileDetails = DistinctFilenames(filenames, &errList)

	// Each file can contain a single element of type 'data', or an array of these elements
	parsedMap = make(ParsedMap)
	for _, file = range fileDetails {
		if b, err = ioutil.ReadFile(file.Name); err != nil {
			errList = append(errList, fmt.Sprintf("reading file: %v", err))
			continue
		}
		config = nil

		// Parse config file as map[string]interface{}, or slice of these
		err = json.Unmarshal(b, &config)
		if err != nil {
			golog.Log.Debugf("Parsing issue, skipping [%s]", file.Name)
			errList = append(errList, fmt.Sprintf("%v, skipping [%s]", err, file.Name))
			continue
		}

		k = reflect.TypeOf(config).Kind().String()
		if k == "map" {
			golog.Log.Debugf("Parsing single element [%s]", file.Name)

			parsed := Parsed{
				FileName:     file.Name,
				DistinctName: file.DistinctName,
				ElementMap:   config.(map[string]interface{}),
			}

			// Find element Id by json tag or field name
			v = nil
			if tag {
				v, _ = parsed.ElementMap[idTag]
			}
			if v == nil {
				v, _ = parsed.ElementMap[idName]
			}
			if v == nil {
				errList = append(errList, fmt.Sprintf("required id parameter %s not found, skipping [%s]",
					idName, file.Name))
				continue
			}
			elementId = fmt.Sprintf("%v", v)

			// Add parsed to array and store in resultMap
			parsedArr, ok = parsedMap[elementId]
			if ok {
				parsedArr = append(parsedArr, parsed)
			} else {
				parsedArr = []Parsed{parsed}
			}
			parsedMap[elementId] = parsedArr

		} else if k == "slice" {
			golog.Log.Debugf("Parsing %d elements [%s]", len(config.([]interface{})), file.Name)

			// Config file contains an array of elements of type 'data'
			for i, v = range config.([]interface{}) {

				parsed := Parsed{
					FileName:     file.Name,
					DistinctName: file.DistinctName,
					Position:     i + 1,
					ElementMap:   v.(map[string]interface{}),
				}

				// Find element Id by json tag or field name
				v = nil
				if tag {
					v, _ = parsed.ElementMap[idTag]
				}
				if v == nil {
					v, _ = parsed.ElementMap[idName]
				}
				if v == nil {
					errList = append(errList, fmt.Sprintf("required id parameter %s not found, skipping [%s:elem#%d]",
						idName, file.Name, parsed.Position))
					continue
				}
				elementId = fmt.Sprintf("%v", v)

				// Add parsed to array and store in resultMap
				parsedArr, ok = parsedMap[elementId]
				if ok {
					parsedArr = append(parsedArr, parsed)
				} else {
					parsedArr = []Parsed{parsed}
				}
				parsedMap[elementId] = parsedArr
			}
		} else {
			errList = append(errList, fmt.Sprintf("parsing config: unrecognized JSON type %q [%s]", k, file.Name))
			continue
		}
	}

	// check for conflicting values and unused parameters
	validateParameters(st, parsedMap, &errList)

	// collapse each element into a single data object and load into resultMap
	resultMap = parseConfig(st, sv, parsedMap, &errList)

	// Compile error list into an error message
	err = errList.Get()

	golog.Log.Debugf("Parsed %d distinct configurations", len(parsedMap))
	return
}

// Parse config dataMap entries corresponding to data object (st, sv) fields
// - Allows values to be numbers, strings, dates, datetimes, or durations
// - Numbers can be specified as a string or value
// - Errors if extra parameters configured
func parseConfig(st reflect.Type, sv reflect.Value, parsedMap ParsedMap, errList *ErrList) (resultMap ResultMap) {
	var err error
	var elementId, filename, tagName, paramValue, paramType string
	var v interface{}
	var parsedArr []Parsed
	var dur time.Duration
	var date time.Time
	var f float64
	var n int64
	var i int
	var ok, clear bool

	// Map for json param names to tag names
	param2tag := make(map[string]string)
	for i = 0; i < st.NumField(); i++ {
		param := st.Field(i)
		tagName, ok = param.Tag.Lookup("json")
		if ok {
			param2tag[param.Name] = tagName
		}
	}

	// If multiple results, we will have to clear 'data' object each iteration
	clear = len(parsedMap) > 1

	resultMap = make(ResultMap)
	for elementId, parsedArr = range parsedMap {

		// Make map of parameter names that will need to be reset after parsing element
		clearParamMap := make(map[string]bool)

		for _, parsed := range parsedArr {
			if parsed.Position == 0 {
				filename = parsed.DistinctName
			} else {
				filename = fmt.Sprintf("%s:elem#%d", parsed.DistinctName, parsed.Position)
			}

			// Iterate through element data fields, parse into correct type
			for i = 0; i < st.NumField(); i++ {
				param := st.Field(i)

				// lookup in ElementMap by tag name first, then by param name
				tagName, ok = param2tag[param.Name]
				if ok {
					v, ok = parsed.ElementMap[tagName]
				}
				if !ok {
					v, ok = parsed.ElementMap[param.Name]
				}

				if ok {
					paramValue = fmt.Sprintf("%v", v)
					paramType = param.Type.Name()
					clearParamMap[param.Name] = true

					switch paramType {
					case "string":
						sv.Field(i).SetString(paramValue)
					case "float64":
						f, err = strconv.ParseFloat(paramValue, 64)
						if err != nil {
							errList.Addf("setting for %s invalid, parameter %s: float %s [%s]",
								elementId, param.Name, paramValue, filename)
							continue
						}
						sv.Field(i).SetFloat(f)
					case "int", "int64":
						n, err = strconv.ParseInt(paramValue, 10, 64)
						if err != nil {
							errList.Addf("setting for %s invalid, parameter %s: integer %s [%s]",
								elementId, param.Name, paramValue, filename)
							continue
						}
						sv.Field(i).SetInt(n)
					case "bool":
						ok, err = strconv.ParseBool(paramValue)
						if err != nil {
							errList.Addf("setting for %s invalid, parameter %s: boolean %s [%s]",
								elementId, param.Name, paramValue, filename)
							continue
						}
						sv.Field(i).SetBool(ok)
					case "Duration":
						dur, err = time.ParseDuration(paramValue)
						if err != nil {
							errList.Addf("setting for %s invalid, parameter %s: duration %s [%s]",
								elementId, param.Name, paramValue, filename)
							continue
						}
						sv.Field(i).Set(reflect.ValueOf(dur))
					case "Time":
						date, err = time.Parse("2006-01-02T15:04:05Z", paramValue)
						if err != nil {
							date, err = time.Parse("2006-01-02", paramValue)
						}
						if err != nil {
							errList.Addf("setting for %s invalid, parameter %s: date %s [%s]",
								elementId, param.Name, paramValue, filename)
							continue
						}
						sv.Field(i).Set(reflect.ValueOf(date))
					default:
						errList.Addf("setting for %s invalid, parameter %s: unsupported type %s [%s]",
							elementId, param.Name, paramType, filename)
						continue
					}
				}
			}
		}
		// Store data object in resultMap
		resultMap[elementId] = sv.Interface()

		// Clear data object for next element Id
		if clear {
			clearConfig(st, sv, clearParamMap)
		}
	}
	return
}

func clearConfig(st reflect.Type, sv reflect.Value, clearParamMap map[string]bool) (err error) {
	var paramType string
	var zeroD time.Duration
	var zeroT time.Time
	var i int
	var ok bool

	// Iterate through data fields, clear settings
	for i = 0; i < st.NumField(); i++ {
		param := st.Field(i)

		ok = clearParamMap[param.Name]
		if ok {
			paramType = param.Type.Name()

			switch paramType {
			case "string":
				sv.Field(i).SetString("")
			case "float64":
				sv.Field(i).SetFloat(0)
			case "int", "int64":
				sv.Field(i).SetInt(0)
			case "bool":
				sv.Field(i).SetBool(false)
			case "Duration":
				sv.Field(i).Set(reflect.ValueOf(zeroD))
			case "Time":
				sv.Field(i).Set(reflect.ValueOf(zeroT))
			default:
				err = fmt.Errorf("unsupported type %s [%s]", paramType, param.Name)
				return
			}
		}
	}

	return
}

// Check data object (st) fields for any conflicting result map values
func validateParameters(st reflect.Type, parsedMap ParsedMap, errList *ErrList) {
	var filenames []string
	var elementId, key, filename, tagName, paramName, paramValue string
	var parsedArr []Parsed
	var paramValuesMap map[string][]string
	var v interface{}
	var i int
	var ok bool

	// Maps for json tag names <-> param names, and valid param names
	tag2param := make(map[string]string)
	param2tag := make(map[string]string)
	validparam := make(map[string]bool)
	for i = 0; i < st.NumField(); i++ {
		param := st.Field(i)
		tagName, ok = param.Tag.Lookup("json")
		if ok {
			tag2param[tagName] = param.Name
			param2tag[param.Name] = tagName
		}
		validparam[param.Name] = true
	}

	// Validate parameters for each element
	for elementId, parsedArr = range parsedMap {

		// make list of values for each parameter, with filenames found in
		elementParamValuesMap := make(map[string]map[string][]string)

		// make map of parameter names, with filenames found in, to see what isn't used
		unusedParamMap := make(map[string][]string)

		// check across all files parsed
		for _, parsed := range parsedArr {
			if parsed.Position == 0 {
				filename = parsed.DistinctName
			} else {
				filename = fmt.Sprintf("%s:elem#%d", parsed.DistinctName, parsed.Position)
			}

			// load elementParamValuesMap to identify possible conflicting values
			for i = 0; i < st.NumField(); i++ {
				param := st.Field(i)
				// lookup in config dataMap by tag name first, then by param name
				tagName, ok = param2tag[param.Name]
				if ok {
					v, ok = parsed.ElementMap[tagName]
				}
				if !ok {
					v, ok = parsed.ElementMap[param.Name]
				}
				if ok {
					paramValue = fmt.Sprintf("%v", v)

					// make list all filenames for each parameter value
					paramValuesMap, ok = elementParamValuesMap[param.Name]
					if !ok {
						paramValuesMap = make(map[string][]string)
						filenames = []string{filename}
					} else {
						filenames, ok = paramValuesMap[paramValue]
						if ok {
							filenames = append(filenames, filename)
						} else {
							filenames = []string{filename}
						}
					}
					paramValuesMap[paramValue] = filenames
					elementParamValuesMap[param.Name] = paramValuesMap
				}
			}

			// load unusedParamMap to identify possible unused parameters
			for key, _ = range parsed.ElementMap {
				// lookup element key by tag name first
				paramName, ok = tag2param[key]
				if !ok {
					paramName = key
				}
				if !validparam[paramName] {
					filenames, ok = unusedParamMap[paramName]
					if ok {
						filenames = append(filenames, filename)
					} else {
						filenames = []string{filename}
					}
					unusedParamMap[paramName] = filenames
				}
			}
		}

		// List errors for conflicting values
		for paramName, paramValuesMap = range elementParamValuesMap {
			// if there are more than one value, it means settings conflict
			if len(paramValuesMap) > 1 {
				var conflicts []string
				for paramValue, filenames = range paramValuesMap {
					conflicts = append(conflicts, fmt.Sprintf("%q [%s]",
						paramValue, strings.Join(filenames, ",")))
				}
				errList.Addf("settings for %s conflict, parameter %s: %s",
					elementId, paramName, strings.Join(conflicts, " != "))
			}
		}

		// List  errors for unused parameters
		for paramName, filenames = range unusedParamMap {
			if len(filenames) == 1 {
				errList.Addf("unused setting for %s, parameter %s [%s]", elementId, paramName, filenames[0])
			} else {
				errList.Addf("unused settings for %s, parameter %s: %d occurences [%s]",
					elementId, paramName, len(filenames), strings.Join(filenames, ","))
			}
		}
	}
	return
}
