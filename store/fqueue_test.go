package store

import "testing"

func TestNewFileQueue(t *testing.T) {
	s := NewFileQueue("/tmp/flemq")
	if s == nil {
		t.Fatalf("NewFileQueue returned nil")
	}
}
