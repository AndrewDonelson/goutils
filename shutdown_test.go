package goutils

import "testing"

// func TestShutdown(t *testing.T) {
// 	DelayShutdown = 5
// 	go CatchShutdown()
// }

func TestShutdownReason(t *testing.T) {
	DelayReason = "Just Because!"
	go CatchShutdown()
}

// func TestShutdownBusy(t *testing.T) {
// 	go CatchShutdown()
// 	go count()
// 	os.Exit(0)
// }
