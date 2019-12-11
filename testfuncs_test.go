package goutils

import (
	"testing"
)

func TestEquals(t *testing.T) {
	a := 1
	Equals(t, a, a)
}

func TestAssert(t *testing.T) {
	a := 1
	Assert(t, (a == 1), "False")
}

func TestOk(t *testing.T) {
	//err := fmt.Errorf("Some Error")
	Ok(t, nil)
}
