package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alainrk/flemq/config"
	"github.com/alainrk/flemq/server"
	"github.com/joho/godotenv"
)

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
	// TODO: Enable on dev only
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config.NewConfig()
	fmt.Println(config)

	server, closer := server.NewServer(config)

	// If killed by SIGTERM
	handleSignals(closer)
	// If died natually
	defer closer()

	server.Run()
}
