package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	handleSignals()

	fmt.Println("starting server")

}

// handleSignals registers signal handlers for shutdown
func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
}
