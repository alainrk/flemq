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

// Receive everything from the client until connection gets closed
// func handleClient(conn net.Conn) {
// 	defer conn.Close()

// 	// Set timeout on r/w operations
// 	conn.SetDeadline(time.Now().Add(RW_TIMEOUT))

// 	buf, err := io.ReadAll(conn)
// 	if err != nil && len(buf) == 0 {
// 		fmt.Println("Error reading from connection:", err.Error())
// 		return
// 	}

// 	fmt.Printf("Received: \"%s\"\n", string(buf))

// 	// Write data back to the client
// 	_, err = conn.Write(buf)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// }

// It works but gets stuck on write side if larger than X bytes
// func handleClient(conn net.Conn) {
//  defer conn.Close()
// 	chunk := make([]byte, RECV_CHUNK_SIZE)

// 	for {
// 		len, err := conn.Read(chunk)

// 		if err != nil {
// 			fmt.Println("Error reading:", err.Error())
// 			break
// 		}

// 		s := string(chunk[:len])

// 		fmt.Println("Stuff", s)
// 		fmt.Println("len", binary.Size(chunk))
// 	}
// }

func handleClient(conn net.Conn) {
	defer conn.Close()

	var received int
	// The buffer grows as we write into it.
	// Ref: https://pkg.go.dev/bytes#Buffer
	buffer := bytes.NewBuffer(nil)

	// Read the data in chunks.
	for {
		// e.g. 8192
		chunk := make([]byte, RECV_CHUNK_SIZE)
		read, err := conn.Read(chunk)
		if err != nil {
			fmt.Println(received, buffer.Bytes(), err)
			return
		}
		received += read
		buffer.Write(chunk[:read])

		if read == 0 || read < RECV_CHUNK_SIZE {
			break
		}
	}

	conn.Write([]byte("OK\n"))
	fmt.Println("Recvd:", received, buffer.String())
}

// func handleClient2(conn net.Conn) {
// 	defer conn.Close()

// 	// Wrap connection with buffered reader so size gets handled automatically
// 	r := bufio.NewReader(conn)

// 	for {
// 		// s, err := r.ReadString('\n')
// 		s, err := r.ReadBytes('\n')
// 		if err != nil {
// 			fmt.Println("Error reading:", err.Error())
// 			break
// 		}
// 		fmt.Println("received:", s)
// 	}
// }

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

	// If killed by SIGTERM
	handleSignals(closer)
	// If died natually
	defer closer()

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
