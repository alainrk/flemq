package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var SERV_ADDR = "localhost:22123"

func consumer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("consumer: stop consuming and exiting...")
			return
		}
	}
}

func producer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("consumer: stop producing and exiting...")
			return
		default:
			conn, err := net.Dial("tcp", SERV_ADDR)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			conn.Write([]byte(`hello world`))
			conn.Close()
			time.Sleep(1 * time.Second)
		}
	}

	// syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}

// handleSignals registers signal handlers for shutdown
func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	handleSignals(cancel)
	wg.Add(2)
	go func() {
		defer wg.Done()
		consumer(ctx)
	}()
	go func() {
		defer wg.Done()
		producer(ctx)
	}()

	wg.Wait()
	fmt.Println("received SIGTERM, exiting")
	os.Exit(0)
}
