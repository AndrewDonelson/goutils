package goutils

import (
	"testing"
)

func TestEquals(t *testing.T) {
	a := 1
	equals(t, a, a)
}

func TestAssert(t *testing.T) {
	a := 1
	assert(t, (a == 1), "False")
}

func TestOk(t *testing.T) {
	//err := fmt.Errorf("Some Error")
	ok(t, nil)
}
