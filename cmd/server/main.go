package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
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
		fmt.Println("received SIGTERM, exiting")
		closer()
		os.Exit(0)
	}()
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New client:", conn.RemoteAddr())
	conn.SetDeadline(time.Now().Add(RW_TIMEOUT))

	var received int
	// The buffer grows as we write into it.
	// Ref: https://pkg.go.dev/bytes#Buffer
	buffer := bytes.NewBuffer(nil)

	// Read the data in chunks.
	for {
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

func main() {
	var err error
	var listener net.Listener

	if TLS_ENABLED {
		cert, err := tls.LoadX509KeyPair(TLS_CERT_FILE, TLS_KEY_FILE)
		if err != nil {
			fmt.Println("Error loading cert:", err)
			return
		}

		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", ADDR, config)
	} else {
		listener, err = net.Listen("tcp", ADDR)
	}

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
