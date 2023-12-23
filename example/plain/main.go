package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var SERV_ADDR = "localhost:22123"

func producer() {

	for {
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

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		producer()
	}()

	wg.Wait()
}
