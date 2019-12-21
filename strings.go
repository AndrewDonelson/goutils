package goutils

import (
	"fmt"
	"regexp"
)

// StringToAlphaNumeric given a string will strip all characters except AlphaNumeric
// and return the result
func StringToAlphaNumeric(s string) (str string, err error) {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	str = reg.ReplaceAllString(s, "")
	return str, nil
}

// StringArrayContains helper function to return true or false is a given string exists in the provided string array
func StringArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// ValidateEmailAddress tests given email for valid email format
func ValidateEmailAddress(email string) (err error) {
	// Validate Email
	minLengthOk := len(email) >= 4

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !minLengthOk || (re.MatchString(email) == false) {
		err = fmt.Errorf("Email is not in a valid format")
		return
	}
	return
}
