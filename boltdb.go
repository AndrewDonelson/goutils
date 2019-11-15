package goutils

import (
	"fmt"
	"os"

	"github.com/AndrewDonelson/golog"
	"github.com/boltdb/bolt"
)

var BoltDB *bolt.DB

func ConnectBolt(file string) (err error) {

	golog.Log.Infof("Connecting to BoltDB at %s", file)
	BoltDB, err = bolt.Open(file, 0644, nil)

	return
}

func InitializeBolt(file string) (err error) {
	golog.Log.Infof("Initializing BoltDB at %s", file)

	// Remove previous DB if exists
	err = os.Remove(file)
	if err != nil && !os.IsNotExist(err) {
		err = fmt.Errorf("deleting previous Bolt DB: %v", err)
		return
	}

	// Create empty file for new db
	_, err = os.Create(file)
	if err != nil {
		err = fmt.Errorf("creating Bolt DB: %v", err)
	}
	return
}
