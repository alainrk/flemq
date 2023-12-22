package main

import (
	"bytes"
	"fmt"
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
func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Set timeout on r/w operations
	conn.SetDeadline(time.Now().Add(RW_TIMEOUT))

	// Read and process data from the client looping through a buffer of 1024 bytes
	var recvd int
	buf := bytes.NewBuffer(nil)

	for {
		chunk := make([]byte, RECV_CHUNK_SIZE)
		n, err := conn.Read(chunk)

		if err != nil {
			fmt.Printf("Err reading chunk: \"%s\"", string(chunk[:n]))
			return
		}

		recvd += n
		buf.Write(chunk[:n])

		if recvd == 0 || recvd < RECV_CHUNK_SIZE {
			break
		}
	}

	fmt.Printf("Received: \"%s\"\n", buf.String())

	// Write data back to the client
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func main() {
	handleSignals()

	// Start TCP server
	listener, err := net.Listen("tcp", ADDR)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

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
