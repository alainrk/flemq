package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

type ClientStatus string

type Server struct {
	clients  map[uuid.UUID]*Client
	listener net.Listener
}

type Client struct {
	Id         uuid.UUID
	Connection net.Conn
}

func NewServer() (server *Server, closer func()) {
	var listener net.Listener
	var err error

	if TLS_ENABLED {
		cert, err := tls.LoadX509KeyPair(TLS_CERT_FILE, TLS_KEY_FILE)
		if err != nil {
			log.Fatalln("Error loading cert:", err)
		}

		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", ADDR, config)
	} else {
		listener, err = net.Listen("tcp", ADDR)
	}

	log.Println("Server is listening on", ADDR)

	if err != nil {
		log.Fatalln("Error:", err)
	}

	closer = func() {
		listener.Close()
	}

	return &Server{
		clients:  make(map[uuid.UUID]*Client),
		listener: listener,
	}, closer
}

func (s Server) Run() {
	for {
		// Accept incoming connections
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		id := s.NewClient(conn)
		go s.HandleClient(id)
	}
}

func (s Server) NewClient(conn net.Conn) uuid.UUID {
	id := uuid.New()
	c := &Client{
		Id:         id,
		Connection: conn,
	}
	s.clients[id] = c
	return id
}

func (s Server) RemoveClient(id uuid.UUID) {
	c := s.clients[id]
	c.Connection.Close()
	delete(s.clients, id)
}

func (s Server) HandleClient(id uuid.UUID) {
	c := s.clients[id]
	defer s.RemoveClient(id)

	log.Println("New client:", c.Connection.RemoteAddr())
	c.Connection.SetDeadline(time.Now().Add(RW_TIMEOUT))

	var received int
	// The buffer grows as we write into it.
	// Ref: https://pkg.go.dev/bytes#Buffer
	buffer := bytes.NewBuffer(nil)

	// Read the data in chunks.
	for {
		chunk := make([]byte, RECV_CHUNK_SIZE)
		read, err := c.Connection.Read(chunk)
		if err != nil {
			log.Println(received, buffer.Bytes(), err)
			return
		}
		received += read
		buffer.Write(chunk[:read])

		if read == 0 || read < RECV_CHUNK_SIZE {
			break
		}
	}

	c.Connection.Write([]byte("OK\n"))
	log.Println("Recvd:", received, buffer.String())
}
