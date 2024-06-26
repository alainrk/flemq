package server

import (
	"bufio"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"time"

	"github.com/alainrk/flemq/config"
	"github.com/alainrk/flemq/flep"
	"github.com/alainrk/flemq/handlers"
	"github.com/google/uuid"
)

type ClientStatus string

type Server struct {
	handlers handlers.Handlers
	listener net.Listener
	clients  map[uuid.UUID]*Client
	config   config.Config
}

type Client struct {
	Connection net.Conn
	FLEPReader *flep.Reader
	Id         uuid.UUID
}

func NewServer(c config.Config) (server *Server, closer func()) {
	var listener net.Listener
	var err error

	if c.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(c.TLS.CertFile, c.TLS.KeyFile)
		if err != nil {
			log.Fatalln("Error loading cert:", err)
		}

		ctls := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, _ = tls.Listen("tcp", c.Addr, ctls)
	} else {
		listener, err = net.Listen("tcp", c.Addr)
	}

	if err != nil {
		log.Fatalln("Error creating a plaintext listener:", err)
	}

	handlers := handlers.NewHandlers(c)

	closer = func() {
		handlers.Close()
		listener.Close()
	}

	log.Println("Server is listening on", c.Addr)

	return &Server{
		config:   c,
		clients:  make(map[uuid.UUID]*Client),
		listener: listener,
		handlers: handlers,
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
		FLEPReader: flep.NewReader(bufio.NewReader(conn)),
	}
	s.clients[id] = c
	return id
}

func (s Server) RemoveClient(id uuid.UUID) {
	c := s.clients[id]
	c.Connection.Close()
	delete(s.clients, id)
	log.Println("Connected clients:", len(s.clients))
}

func (s Server) HandleClient(id uuid.UUID) {
	c := s.clients[id]
	defer s.RemoveClient(id)

	log.Println("New client:", c.Connection.RemoteAddr())

repl:
	// Read using the flep reader.
	for {
		c.Connection.SetDeadline(time.Now().Add(s.config.Connection.RWTimeout))

		req, err := c.FLEPReader.ReadRequest()
		if err != nil {
			if errors.As(err, &flep.FlepError{}) {
				log.Println("Flep Error:", err)
				fr := flep.SimpleErrorResponse(err.Error())
				c.Connection.Write(fr)
				continue
			}
			log.Println("Connection Error:", err)
			break repl
		}

		// TODO: Remove or set log debug level.
		log.Printf("Handling request: %s", req)

		switch req.Command {

		case flep.CommandPush:
			offset, err := s.handlers.HandlePush(req)
			if err != nil {
				log.Println("Error:", err)
				fr := flep.SimpleErrorResponse(err.Error())
				c.Connection.Write(fr)
				continue
			}
			fr := flep.IntResponse(int64(offset))
			c.Connection.Write(fr)
			// TODO: Remove or set log debug level.
			log.Printf("Written response: %s", string(fr))

		case flep.CommandPick:
			res, err := s.handlers.HandlePick(req)
			if err != nil {
				log.Println("Error:", err)
				fr := flep.SimpleErrorResponse(err.Error())
				c.Connection.Write(fr)
				continue
			}
			fr := flep.SimpleBytesResponse(res)
			c.Connection.Write(fr)

		case flep.CommandSubscribe:
			// Long running command, so we reset the deadline
			// and leave this connection open to be handled.
			c.Connection.SetDeadline(time.Time{})

			done := make(chan bool)

			go s.handlers.HandleSubscribe(c.Connection, req, done)
			go waitDisconnect(c.Connection, done)

			// Wait for any of the two to be done/close/error out.
			<-done

			exit(c.Connection)
			break repl

		case flep.CommandExit:
			exit(c.Connection)
			break repl
		}
	}
}

func waitDisconnect(c net.Conn, done chan bool) {
	_, err := c.Read([]byte{0})
	// TODO: Improvable checking for errors?
	// For subscribers specifically I don't care as I just need to close as soon as there's any error.
	if err != nil {
		done <- true
	}
}

func exit(c net.Conn) {
	log.Println("Client exiting:", c.RemoteAddr())
	fr := flep.SimpleStringResponse("OK")
	c.Write(fr)
}
