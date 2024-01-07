package store

import "io"

type QueueStore interface {
	Write(reader io.Reader) (offset uint64, err error)
	Read(offset uint64, writer io.Writer) error
	Close() error
}
