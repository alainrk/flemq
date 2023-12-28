package server

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/alainrk/flemq/config"
	"github.com/alainrk/flemq/flep"
	"github.com/alainrk/flemq/handlers"
	"github.com/alainrk/flemq/store"
	"github.com/google/uuid"
)

type ClientStatus string

type Server struct {
	config   config.Config
	clients  map[uuid.UUID]*Client
	listener net.Listener
	handlers handlers.Handlers
}

type Client struct {
	Id         uuid.UUID
	Connection net.Conn
	FLEPReader *flep.Reader
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
		listener, err = tls.Listen("tcp", c.Addr, ctls)
		if err != nil {
			log.Fatalln("Error creating a TLS listener:", err)
		}
	} else {
		listener, err = net.Listen("tcp", c.Addr)
		if err != nil {
			log.Fatalln("Error creating a plaintext listener:", err)
		}
	}

	closer = func() {
		listener.Close()
	}

	log.Println("Server is listening on", c.Addr)

	return &Server{
		config:   c,
		clients:  make(map[uuid.UUID]*Client),
		listener: listener,
		handlers: handlers.NewHandlers(store.NewMemoryQueueStore()),
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
	var (
		c *Client = s.clients[id]
	)
	defer s.RemoveClient(id)

	log.Println("New client:", c.Connection.RemoteAddr())

repl:
	// Read using the flep reader.
	for {
		c.Connection.SetDeadline(time.Now().Add(s.config.Connection.RWTimeout))

		req, err := c.FLEPReader.ReadRequest()
		if err != nil {
			if errors.As(err, &flep.FlepError{}) {
				log.Println("Error:", err)
				c.Connection.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err)))
				continue
			}
			log.Println("Error:", err)
			break repl
		}

		switch req.Command {

		case flep.CommandPush:
			offset, err := s.handlers.HandlePush(req)
			if err != nil {
				log.Println("Error:", err)
				c.Connection.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err)))
				continue
			}
			c.Connection.Write([]byte(fmt.Sprintf(":%d\r\n", offset)))

		case flep.CommandPick:
			res, err := s.handlers.HandlePick(req)
			if err != nil {
				log.Println("Error:", err)
				c.Connection.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err)))
				continue
			}
			c.Connection.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(res), res)))

		case flep.CommandSubscribe:
			// Long running command, so we reset the deadline
			// and leave this connection open to be handled.
			c.Connection.SetDeadline(time.Time{})
			err := s.handlers.HandleSubscribe(c.Connection, req)
			if err != nil {
				log.Println("Error:", err)
				c.Connection.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err)))
				continue
			}
			// TODO: Decide how to handle the end of the subscription (if any).
			// c.Connection.Write([]byte("OK")))

		case flep.CommandExit:
			log.Println("Client exiting:", c.Connection.RemoteAddr())
			c.Connection.Write([]byte("+OK\r\n"))
			break repl
		}
	}
}
