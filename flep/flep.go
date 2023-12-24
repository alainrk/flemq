package flep

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type CommandType string

const (
	CommandPush CommandType = "PUSH"
	CommandPick CommandType = "PICK"
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
		return req, fmt.Errorf("invalid command")
	}

	// First arg is the command
	command := CommandType(strings.ToUpper(string(args[0])))

	switch command {
	// PUSH topic message
	case CommandPush:
		if len(args) != 3 {
			return req, fmt.Errorf("invalid PUSH command, must follow: `PUSH topic message`")
		}
		req.Command = command
		req.Args = args[1:]
		return req, nil

	// PICK topic offset
	case CommandPick:
		if len(args) != 3 {
			return req, fmt.Errorf("invalid PICK command, must follow: `PICK topic offset`")
		}
		req.Command = command
		req.Args = args[1:]
		return req, nil
	}

	return req, fmt.Errorf("invalid command")
}
