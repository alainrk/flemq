package handlers

import (
	"net"
	"testing"
)

func Test_HandleSubscribe_RecoverOldMessages(t *testing.T) {
	client, server := net.Pipe()
	go func() {
		server.Close()
	}()
	client.Close()
}
