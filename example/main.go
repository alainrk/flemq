package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func consumer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("consumer: stop consuming and exiting...")
			return
		}
	}
}

func producer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("producer: stop producing and exiting...")
			return
		}
	}
}

// handleSignals registers signal handlers for shutdown
func handleSignals(cancel context.CancelFunc, wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer wg.Done()
		<-c
		cancel()
	}()
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(3)
	handleSignals(cancel, &wg)
	go consumer(ctx, &wg)
	go producer(ctx, &wg)

	wg.Wait()
	fmt.Println("received SIGTERM, exiting")
}
