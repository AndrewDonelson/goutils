package goutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileDetail ...
type FileDetail struct {
	Name           string
	DistinctName   string
	FullPath       string
	PathComponents []string
}

// DistinctFilenames Create Parsed{} struct for each filename, validating and formatting names
func DistinctFilenames(filenames []string, errList *ErrList) (fileDetails []FileDetail) {
	var (
		err                           error
		file                          FileDetail
		items, newitems, remainitems  []string
		filename, fullpath, dir, name string
		i                             int
		ok, distinct                  bool
	)

	// Map to determine distinct names
	usednames := make(map[string][]string)

	// Map to make sure full-path names are unique
	fullnameMap := make(map[string]FileDetail)

	for _, filename = range filenames {

		// Make sure each filename is a valid file
		fullpath, _, err = ValidateFile(filename)
		if err != nil {
			errList.Add(err)
			continue
		}

		// Check for duplicates
		file, ok = fullnameMap[fullpath]
		if ok {
			*errList = append(*errList, fmt.Sprintf("file is duplicate of %s, skipping [%s]", file.Name, filename))
			continue
		}

		// Split base from dir to create initial FileDetail
		dir, name = filepath.Split(fullpath)
		if len(name) == 0 {
			errList.Addf("invalid file [%s]", filename)
			continue
		}
		file = FileDetail{
			Name:         filename,
			DistinctName: name,
			FullPath:     fullpath,
		}

		// Make list of file directory components in absolute path
		for i = 0; i < MaxIterations; i++ {
			dir = filepath.Dir(dir)
			name = filepath.Base(dir)
			file.PathComponents = append(file.PathComponents, name)
			if len(dir) == 1 {
				break
			}
		}
		fullnameMap[fullpath] = file

		// List base names used, to see if names are distinct
		items, ok = usednames[file.DistinctName]
		if ok {
			items = append(items, fullpath)
		} else {
			items = []string{fullpath}
		}
		usednames[file.DistinctName] = items
	}

	for i = 0; i < MaxIterations; i++ {
		distinct = true

		// Generate distinct names across all files
		for name, items = range usednames {

			// If multiple items, we don't have distinct filenames
			if len(items) > 1 {
				distinct = false
				for _, fullpath = range items {
					file, ok = fullnameMap[fullpath]
					if !ok {
						fmt.Printf("internal error processing file %s\n", fullpath)
						continue
					}
					if len(file.PathComponents) > 0 {
						// prepend distinct name with first path component from list
						file.DistinctName = filepath.Join(file.PathComponents[0], file.DistinctName)
						file.PathComponents = file.PathComponents[1:]

						// store new distinct name
						newitems, ok = usednames[file.DistinctName]
						if ok {
							newitems = append(newitems, fullpath)
						} else {
							newitems = []string{fullpath}
						}
						usednames[file.DistinctName] = newitems

						// Update results
						fullnameMap[fullpath] = file
					} else {
						// exhausted path components, so distinct name remains unchanged
						remainitems = append(remainitems, fullpath)
					}
				}

				if len(remainitems) > 0 {
					usednames[name] = remainitems
				} else {
					delete(usednames, name)
				}
			}
		}

		// Having checked across all files, each name used only once
		if distinct {
			break
		}
	}

	// Copy results from fullnameMap
	for _, file = range fullnameMap {
		fileDetails = append(fileDetails, file)
	}
	return
}

// ValidateFile Validate if exists, and is type file
func ValidateFile(filename string) (fullpath string, fileInfo os.FileInfo, err error) {

	// Is filename valid
	if len(filename) == 0 {
		err = fmt.Errorf("{file} required")
		return
	}
	fullpath, err = filepath.Abs(filename)
	if err != nil {
		err = fmt.Errorf("%v [%s]", err, filename)
		return
	}

	// Is accessible and type file
	fileInfo, err = os.Stat(fullpath)
	if os.IsNotExist(err) {
		err = fmt.Errorf("invalid, no such file [%s]", filename)
		return
	} else if os.IsPermission(err) {
		err = fmt.Errorf("invalid, permission denied [%s]", filename)
		return
	} else if err != nil {
		err = fmt.Errorf("invalid, %v", err)
		return
	} else if fileInfo.IsDir() {
		err = fmt.Errorf("invalid, file is a directory [%s]", filename)
	}
	return
}

// ValidateDir Validate if exists, and is type directory
func ValidateDir(filename string) (fullpath string, dirInfo os.FileInfo, err error) {

	// Is filename valid
	if len(filename) == 0 {
		err = fmt.Errorf("{directory} required")
		return
	}
	fullpath, err = filepath.Abs(filename)
	if err != nil {
		err = fmt.Errorf("%v [%s]", err, filename)
		return
	}

	// Is accessible and type file
	dirInfo, err = os.Stat(fullpath)
	if os.IsNotExist(err) {
		err = fmt.Errorf("invalid, no such directory [%s]", filename)
		return
	} else if os.IsPermission(err) {
		err = fmt.Errorf("invalid, permission denied [%s]", filename)
		return
	} else if err != nil {
		err = fmt.Errorf("invalid, %v", err)
		return
	} else if !dirInfo.IsDir() {
		err = fmt.Errorf("invalid, is not a directory [%s]", filename)
	}
	return
}

// ValidateFileOrParentDir Validate if exists, and is type file
// Or parent directory valid, because file is to be created
func ValidateFileOrParentDir(filename string) (fullpath string, err error) {
	var fileInfo os.FileInfo
	var dirname string

	// Is filename valid
	if len(filename) == 0 {
		err = fmt.Errorf("{file} required")
		return
	}
	fullpath, err = filepath.Abs(filename)
	if err != nil {
		err = fmt.Errorf("%v [%s]", err, filename)
		return
	}

	// Is parent directory valid
	dirname = filepath.Dir(fullpath)
	_, err = os.Stat(dirname)
	if os.IsNotExist(err) {
		err = fmt.Errorf("parent directory %s does not exist", dirname)
		return
	}

	// If exists, is accessible and type file
	fileInfo, err = os.Stat(fullpath)
	if os.IsNotExist(err) {
		err = nil
		return
	} else if os.IsPermission(err) {
		err = fmt.Errorf("invalid, permission denied [%s]", filename)
		return
	} else if err != nil {
		err = fmt.Errorf("invalid, %v", err)
		return
	} else if fileInfo.IsDir() {
		err = fmt.Errorf("invalid, file is a directory [%s]", filename)
	}
	return
}
