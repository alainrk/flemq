package store

import (
	"fmt"
	"io"
	"sync"
)

type QueueStore interface {
	Write(reader io.Reader) (offset uint64, err error)
	Read(offset uint64, writer io.Writer) error
}

type MemoryQueueStore struct {
	mu sync.RWMutex

	counter uint64
	data    map[uint64][]byte
}

func NewMemoryQueueStore() *MemoryQueueStore {
	return &MemoryQueueStore{
		data: make(map[uint64][]byte),
	}
}

func (s *MemoryQueueStore) Write(reader io.Reader) (offset uint64, err error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return 0, err
	}

	// Restrict critical section to the minimum
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[s.counter] = buf
	s.counter++

	return s.counter - 1, nil
}

func (s *MemoryQueueStore) Read(offset uint64, writer io.Writer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[offset]; !ok {
		return fmt.Errorf("offset %d not found", offset)
	}

	_, err := writer.Write(s.data[offset])
	return err
}
