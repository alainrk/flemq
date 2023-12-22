package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ADDR = ":22123"
var RW_TIMEOUT = 60 * time.Second
var RECV_BUF_SIZE = 4096
var RECV_CHUNK_SIZE = 32

// handleSignals registers signal handlers for shutdown
func handleSignals(closer func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		closer()
		os.Exit(0)
	}()
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Set timeout on r/w operations
	conn.SetDeadline(time.Now().Add(RW_TIMEOUT))

	buf, err := io.ReadAll(conn)
	if err != nil && len(buf) == 0 {
		fmt.Println("Error reading from connection:", err.Error())
		return
	}

	fmt.Printf("Received: \"%s\"\n", string(buf))

	// Write data back to the client
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func main() {
	// Start TCP server
	listener, err := net.Listen("tcp", ADDR)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	closer := func() {
		fmt.Println("Closing listener...")
		listener.Close()
	}

	handleSignals(closer)

	fmt.Println("Server is listening on", ADDR)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn)
	}
}
