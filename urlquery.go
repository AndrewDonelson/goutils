package goutils

import (
	"fmt"
	"net/http"

	"github.com/AndrewDonelson/golog"
)

// EnsureNoQueryParameters is a help function for endpoints that require no query parameters
func EnsureNoQueryParameters(r *http.Request) (bool, error) {
	size := len(r.URL.RawQuery)
	if size == 0 {
		return true, nil
	}
	return false, fmt.Errorf("Endpoint does not accept any parameters")
}

// GetQueryParameter is a helper function to get query parameters by name and return
// error ONLY if they are required and not present and DO NOT have a defined default value.
//
// Example:
// ```
// qpApple,ok := GetQueryParameter(r,"apple",false,true,"Granny Smith")
// if err != nil {
// 		HandleError(w, http.StatusInternalServerError.RequestURI, err)
// 		return
// }
// ```
func GetQueryParameter(r *http.Request, name string, required bool, defaultOk bool, def string) (string, error) {

	var (
		ok   bool
		keys []string
	)

	// Check the request for the query parameter
	keys, ok = r.URL.Query()[name]

	if !ok {
		keys = append(keys, "")

		// Is the missing query parameter required?
		if required {

			// Can we set a default value for the parameter?
			if !defaultOk {
				return "", fmt.Errorf("Required query parameter [%s] is not present & default value not allowed", name)
			}

			keys[0] = def
			golog.Log.Warningf("Required query parameter [%s] was not present, set to default value [%s]", name, def)
		}

		//param not present but not required, return default or empty
		if keys == nil || len(keys[0]) < 1 {
			// Can we set a default value for the parameter?
			if defaultOk {
				keys[0] = def
				golog.Log.Debugf("Optional query parameter [%s] was not present, set to default value [%s]", name, def)
			} else {
				golog.Log.Debugf("Optional query parameter [%s] was not present, default not allowed", name)
				keys[0] = ""
			}
		}
	} else {
		golog.Log.Debugf("Query parameter [%s] is present with a value of [%v]", name, keys[0])
	}

	// Query()[name] will return an array of items,
	// we only want the single item.
	key := keys[0]

	// check for nil
	if keys[0] == "" {
		golog.Log.Debugf("Returning URL Param [%s] of empty [nil]", name)
	} else {
		golog.Log.Debugf("Returning URL Param [%s] value of [%s]", name, string(key))
	}

	return key, nil
}
