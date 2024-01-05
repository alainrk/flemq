package store

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/alainrk/flemq/common"
)

func TestMemoryQueueStore_WriteAndRead(t *testing.T) {
	store := NewMemoryQueue()

	// Test Write and Read
	data := []byte("test data")
	offset, err := store.Write(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var readBuffer bytes.Buffer
	err = store.Read(offset, &readBuffer)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !bytes.Equal(data, readBuffer.Bytes()) {
		t.Fatalf("Read data does not match written data")
	}
}

func TestMemoryQueueStore_ReadNonExistentOffset(t *testing.T) {
	store := NewMemoryQueue()

	// Test Read with non-existent offset
	nonExistentOffset := uint64(123)
	var readBuffer bytes.Buffer
	err := store.Read(nonExistentOffset, &readBuffer)
	if err == nil {
		t.Fatalf("Expected error for non-existent offset, but got nil")
	}

	if _, ok := err.(common.OffsetNotFoundError); !ok {
		t.Fatalf("Expected error of type OffsetNotFoundError, but got %T", err)
	}
}

func TestMemoryQueueStore_ConcurrentWritesAndReads(t *testing.T) {
	store := NewMemoryQueue()

	// Test concurrent Writes and Reads
	n := 200
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			data := []byte(fmt.Sprintf("test data %d", index))
			offset, err := store.Write(bytes.NewReader(data))
			if err != nil {
				t.Errorf("Write failed: %v", err)
			}

			var readBuffer bytes.Buffer
			err = store.Read(offset, &readBuffer)
			if err != nil {
				t.Errorf("Read failed: %v", err)
			}

			if !bytes.Equal(data, readBuffer.Bytes()) {
				t.Errorf("Read data does not match written data")
			}
		}(i)
	}

	wg.Wait()
}
