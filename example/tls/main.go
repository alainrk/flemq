package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"
	"time"
)

var SERV_ADDR = "localhost:22123"
var TLS_CERT_FILE = "cert/cert.pem"

func producer() {
	cert, err := os.ReadFile(TLS_CERT_FILE)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		fmt.Println("Error: failed to append certificate")
		return
	}
	config := &tls.Config{RootCAs: certPool}

	for {
		conn, err := tls.Dial("tcp", SERV_ADDR, config)
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
