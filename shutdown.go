package goutils

import (
	"github.com/AndrewDonelson/golog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var DelayShutdown sync.WaitGroup
var DelayReason string

// CatchShutdown is a helper function called via goroutine when you start your app
// the only parameter required is the function you want to call to cleanup your app
// example:
// go goutils.CatchShutdown(app.Shutdown())
func CatchShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until SIGINT or SIGTERM is received.
	s := <-c
	golog.Log.Noticef("Gracefully exiting: %v", s)

	// Allow 5 seconds to complete operation if in progress
	go func() {
		time.Sleep(5 * time.Second)
		DelayShutdown.Done()
	}()
	if len(DelayReason) > 0 {
		golog.Log.Infof("Completing: %s", DelayReason)
	}
	DelayShutdown.Wait()

	// Close MySql connections
	// if DbRW != nil {
	// 	DbRW.Close()
	// }
	// if DbRO != nil {
	// 	DbRO.Close()
	// }

	// Close Bolt connections
	// if BoltA != nil {
	// 	BoltA.Close()
	// }
	// if BoltB != nil {
	// 	BoltB.Close()
	// }
	// if BoltC != nil {
	// 	BoltC.Close()
	// }

	// Close Redis connections
	// if RedisPool != nil {
	// 	RedisPool.Empty()
	// }

	golog.Log.Info("Goodbye")
	os.Exit(0)
}
