package goutils

import (
	"fmt"
	"strings"
)

// ErrList Struct to accumulate a list of errors
type ErrList []string

// Add Append error to list
func (errList *ErrList) Add(v ...interface{}) {
	*errList = append(*errList, fmt.Sprint(v...))
}

// Addf Append formatted error to list
func (errList *ErrList) Addf(text string, v ...interface{}) {
	*errList = append(*errList, fmt.Sprintf(text, v...))
}

// Get to Compile list of error strings into an error message
func (errList ErrList) Get() (err error) {
	var i int

	if len(errList) > 0 {
		if len(errList) == 1 {
			err = fmt.Errorf(errList[0])
		} else {
			for i = range errList {
				errList[i] = fmt.Sprintf("(#%d) %s", i+1, errList[i])
			}
			err = fmt.Errorf("multiple errors\n%s", strings.Join(errList, "\n"))
		}
	}
	return
}
