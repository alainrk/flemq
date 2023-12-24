package flep

import (
	"bufio"
	"bytes"
	"fmt"
)

type CommandType string

const (
	CommandPush CommandType = "PUSH"
)

type Request struct {
	Command CommandType
	Args    [][]byte
}

type Reader struct {
	*bufio.Reader
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{r}
}

// ReadRequest reads a command from the reader and returns a valid Request
// or an error. ReadRequest will only reject if the command is invalid from
// a syntax perspective.
func (r *Reader) ReadRequest() (Request, error) {
	var req Request

	l, _, err := r.ReadLine()
	if err != nil {
		return req, err
	}

	// Split by whitespace
	args := bytes.Split(l, []byte(" "))
	if len(args) == 0 {
		return req, fmt.Errorf("Invalid command")
	}

	// First arg is the command
	command := CommandType(args[0])

	switch command {
	case CommandPush:
		if len(args) != 2 {
			return req, fmt.Errorf("Invalid command")
		}
		req.Command = command
		req.Args = [][]byte{args[1]}
		return req, nil
	}

	return req, fmt.Errorf("Invalid command")
}
