package goutils

import "testing"

func TestBoltConnect(t *testing.T) {
	err := ConnectBolt("test/boltdb.data")
	Equals(t, err, nil)
}

func TestInitializeBolt(t *testing.T) {
	err := InitializeBolt("test/boltdb.data")
	Equals(t, err, nil)
}
