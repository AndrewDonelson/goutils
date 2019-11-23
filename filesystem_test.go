package goutils

import "testing"

func TestValidateFile(t *testing.T) {
	// test empty filename
	_, _, err := ValidateFile("")
	equals(t, err.Error(), "{file} required")

	// test file does not exist
	_, _, err = ValidateFile("noexist.txt")
	equals(t, err.Error(), "invalid, no such file [noexist.txt]")

	// test valid file that exists
	_, _, err = ValidateFile("test/readonly.txt")
	ok(t, err)

	// test error file is directory
	_, _, err = ValidateFile("test")
	equals(t, err.Error(), "invalid, file is a directory [test]")
}

func TestValidateDir(t *testing.T) {
	// test empty filename
	_, _, err := ValidateDir("")
	equals(t, err.Error(), "{directory} required")

	// test file does not exist
	_, _, err = ValidateDir("noexist.txt")
	equals(t, err.Error(), "invalid, no such directory [noexist.txt]")

	// test error is not a directory
	_, _, err = ValidateDir("/test/test")
	equals(t, err.Error(), "invalid, no such directory [/test/test]")

	// test valid file that exists
	_, _, err = ValidateDir("test/readonly.txt")
	equals(t, err.Error(), "invalid, is not a directory [test/readonly.txt]")

}
