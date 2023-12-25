package server

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/alainrk/flemq/commands"
	"github.com/alainrk/flemq/config"
	"github.com/alainrk/flemq/flep"
	"github.com/alainrk/flemq/store"
	"github.com/google/uuid"
)

var RW_TIMEOUT = 60 * time.Second
var RECV_CHUNK_SIZE = 1024

type ClientStatus string

type Server struct {
	clients  map[uuid.UUID]*Client
	listener net.Listener
	commands commands.Commands
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
		clients:  make(map[uuid.UUID]*Client),
		listener: listener,
		commands: commands.NewCommands(store.NewMemoryQueueStore()),
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
	// XXX: Experimenting with leaving the client connected until timeout/EXIT cmd.
	// defer s.RemoveClient(id)

	log.Println("New client:", c.Connection.RemoteAddr())
	c.Connection.SetDeadline(time.Now().Add(RW_TIMEOUT))

	// Read using the flep reader.
	for {
		req, err := c.FLEPReader.ReadRequest()
		if err != nil {
			log.Println("Error:", err)
			return
		}

		switch req.Command {
		case flep.CommandPush:
			s.commands.HandlePush(req)
		case flep.CommandPick:
			s.commands.HandlePick(req)
		case flep.CommandExit:
			log.Println("Client exiting:", c.Connection.RemoteAddr())
			s.RemoveClient(id)
			return
		}
	}
}
