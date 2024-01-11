package store

import (
	"bytes"
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

func TestFileQueue(t *testing.T) {
	var (
		err    error
		offset uint64
		buf    bytes.Buffer
	)

	d := [][]byte{
		[]byte(`0. some data xx000`),
		[]byte(`1. some data xx001`),
		[]byte(`2. some data xx002`),
	}

	testFolder := fmt.Sprintf("/tmp/flemq_test_%d", rand.Int())
	s := NewFileQueue(testFolder)
	defer os.RemoveAll(testFolder)

	for i := 0; i < len(d); i++ {
		offset, err = s.Write(bytes.NewReader(d[i]))
		if err != nil {
			t.Fatalf("Error writing %d data, exited with error: %v", i, err)
		}

		if offset != uint64(i) {
			t.Fatalf("Expected offset to be %d, got %d", i, offset)
		}
	}

	for i := 0; i < len(d); i++ {
		buf.Reset()
		err = s.Read(uint64(i), &buf)
		if err != nil {
			t.Fatalf("Error reading %d data, exited with error: %v", i, err)
		}

		if !bytes.Equal(buf.Bytes(), d[i]) {
			t.Fatalf("Test %d, expected %s, got %s", i, d[i], buf.Bytes())
		}
	}
}

// TODO: Test persistence of data at restart
func TestFileQueuePersistence(t *testing.T) {
	var (
		err    error
		offset uint64
		buf    bytes.Buffer
	)

	d := [][]byte{
		[]byte(`0. some data xx000`),
		[]byte(`1. some data xx001`),
		[]byte(`2. some data xx002`),
	}

	testFolder := fmt.Sprintf("/tmp/flemq_test_%d", rand.Int())
	defer os.RemoveAll(testFolder)

	// First queue creation
	s := NewFileQueue(testFolder)

	for i := 0; i < len(d); i++ {
		offset, err = s.Write(bytes.NewReader(d[i]))
		if err != nil {
			t.Fatalf("Error writing %d data, exited with error: %v", i, err)
		}

		if offset != uint64(i) {
			t.Fatalf("Expected offset to be %d, got %d", i, offset)
		}
	}

	s.Close()

	// Second queue creation
	s = NewFileQueue(testFolder)

	for i := 0; i < len(d); i++ {
		buf.Reset()
		err = s.Read(uint64(i), &buf)
		if err != nil {
			t.Fatalf("Error reading %d data, exited with error: %v", i, err)
		}

		if !bytes.Equal(buf.Bytes(), d[i]) {
			t.Fatalf("Test %d, expected %s, got %s", i, d[i], buf.Bytes())
		}
	}
}

// TODO: Test for non existing offset
