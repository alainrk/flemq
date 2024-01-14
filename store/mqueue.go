package store

import (
	"fmt"
	"io"
	"sync"

	"github.com/alainrk/flemq/common"
)

type MemoryQueue struct {
	mu sync.RWMutex

	counter uint64
	data    map[uint64][]byte
}

// NewMemoryQueue creates a new memory queue.
func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		data: make(map[uint64][]byte),
	}
}

// Write writes the data from the reader into the queue.
func (s *MemoryQueue) Write(reader io.Reader) (offset uint64, err error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return 0, err
	}

	// TODO: What about an actor model here instead?
	// Restrict critical section to the minimum
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[s.counter] = buf
	s.counter++

	return s.counter - 1, nil
}

// Read reads the data at the given offset into the writer.
func (s *MemoryQueue) Read(offset uint64, writer io.Writer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data[offset]; !ok {
		return common.OffsetNotFoundError{Err: fmt.Errorf("offset %d not found", offset)}
	}

	_, err := writer.Write(s.data[offset])
	return err
}

// Close closes the queue.
func (s *MemoryQueue) Close() error {
	return nil
}
