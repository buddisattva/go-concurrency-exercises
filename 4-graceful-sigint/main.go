//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	exitChan := make(chan int)

	// Run the process (blocking)
	go proc.Run()

	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				go proc.Stop()

				for {
					s := <-signalChan
					switch s {
					case syscall.SIGINT:
						exitChan <- 0
					}
				}
			}
		}
	}()

	code := <-exitChan
	fmt.Println("\nSIGINT is called again, just kill the program... (last resort)")
	os.Exit(code)
}
