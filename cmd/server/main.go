package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ADDR = ":22123"
var RW_TIMEOUT = 60 * time.Second
var RECV_CHUNK_SIZE = 1024

var TLS_ENABLED = true
var TLS_CERT_FILE = "cert/cert.pem"
var TLS_KEY_FILE = "cert/key.pem"

// handleSignals registers signal handlers for shutdown
func handleSignals(closer func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("received SIGTERM, exiting")
		closer()
		os.Exit(0)
	}()
}

func main() {
	server, closer := NewServer()

	// If killed by SIGTERM
	handleSignals(closer)
	// If died natually
	defer closer()

	server.Run()
}
