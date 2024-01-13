package main

import (
	"flag"
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

// setup loads environment variables and flags
func setup() {
	devMode := flag.Bool("dev", false, "Enable development mode (loads .env file)")
	tlsEnabled := flag.Bool("tls", false, "Enable TLS (overrides env var)")
	flag.Parse()

	if *devMode {
		log.Println("Development mode enabled, loading .env file...")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if *tlsEnabled {
		os.Setenv("FLEMQ_TLS_ENABLED", "true")
	}
}

func main() {
	// Parse flags and load environment variables
	setup()

	config := config.NewConfig()
	server, closer := server.NewServer(config)

	// If killed by SIGTERM
	handleSignals(closer)
	// If died naturally
	defer closer()

	server.Run()
}
