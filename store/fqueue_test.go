package store

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

// Write(reader io.Reader) (offset uint64, err error)
// Read(offset uint64, writer io.Writer) error
// Close() error

func TestNewFileQueue(t *testing.T) {
	var err error

	// Generate a random folder name
	testFolder := fmt.Sprintf("/tmp/flemq_test_%d", rand.Int())
	defer os.RemoveAll(testFolder)

	s := NewFileQueue(testFolder)
	if s == nil {
		t.Fatalf("NewFileQueue returned nil")
	}

	dataFile := testFolder + "/data"
	_, err = os.Stat(dataFile)
	if err != nil {
		t.Fatalf("Error checking data file %s, exited with error: %v", testFolder, err)
	}

	indexFile := testFolder + "/index"
	_, err = os.Stat(indexFile)
	if err != nil {
		t.Fatalf("Error checking index file %s, exited with error: %v", testFolder, err)
	}

	if s.offset != 0 {
		t.Fatalf("Expected offset to be 0, got %d", s.offset)
	}
}

func TestOffsetAtStartup(t *testing.T) {
	var err error

	// Generate a random folder name
	testFolder := fmt.Sprintf("/tmp/flemq_test_%d", rand.Int())
	defer os.RemoveAll(testFolder)

	s := NewFileQueue(testFolder)
	if s == nil {
		t.Fatalf("NewFileQueue returned nil")
	}

	offset, err := s.getOffsetAtStartup()
	if err != nil {
		t.Fatalf("Error getting offset at startup, exited with error: %v", err)
	}

	if offset != 0 {
		t.Fatalf("Expected offset to be 0, got %d", offset)
	}
}
