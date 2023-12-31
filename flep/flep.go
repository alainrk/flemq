package flep

import (
	"bufio"
	"bytes"
	"strings"
)

type CommandType string

const (
	CommandPush      CommandType = "PUSH"
	CommandPick      CommandType = "PICK"
	CommandSubscribe CommandType = "SUBSCRIBE"
	CommandExit      CommandType = "EXIT"
)

type Request struct {
	Command CommandType
	Args    [][]byte
}

type Reader struct {
	*bufio.Reader
}

type FlepError struct {
	message string
}

func NewFlepError(message string) FlepError {
	return FlepError{
		message: message,
	}
}

func (e FlepError) Error() string {
	return e.message
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
		return req, NewFlepError("invalid command")
	}

	// First arg is the command
	command := CommandType(strings.ToUpper(string(args[0])))

	switch command {
	// PUSH topic message
	case CommandPush:
		if len(args) != 3 {
			return req, NewFlepError("invalid PUSH command, must follow: `PUSH topic message`")
		}
		req.Command = command
		req.Args = args[1:]
		return req, nil

	// PICK topic offset
	case CommandPick:
		if len(args) != 3 {
			return req, NewFlepError("invalid PICK command, must follow: `PICK topic offset`")
		}
		req.Command = command
		req.Args = args[1:]
		return req, nil

	// SUBSCRIBE topic [from_offset_included=0]
	case CommandSubscribe:
		// Default from_offset_included to 0
		if len(args) == 2 {
			args = append(args, []byte("0"))
		}
		if len(args) != 3 {
			return req, NewFlepError("invalid SUBSCRIBE command, must follow: `SUBSCRIBE topic from_offset_included`")
		}
		req.Command = command
		req.Args = args[1:]
		return req, nil

	// EXIT
	case CommandExit:
		if len(args) != 1 {
			return req, NewFlepError("invalid EXIT command, must follow: `EXIT`")
		}
		req.Command = command
		return req, nil
	}

	return req, NewFlepError("invalid command")
}
