package flep

import (
	"bufio"
	"fmt"
)

type CommandType string

const (
	CommandX CommandType = "X"
	CommandY CommandType = "Y"
)

type Request struct {
	Command CommandType
	Args    []string
}

type Reader struct {
	*bufio.Reader
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{r}
}

func (r *Reader) ReadCommand() ([]byte, error) {
	l, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	switch l[0] {
	case '+':
		return []byte("simple string"), nil
	case ':':
		return []byte("integer given"), nil
	default:
		return nil, fmt.Errorf("unknown command: %s", l)
	}
}
