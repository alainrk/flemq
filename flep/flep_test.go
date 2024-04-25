package flep

import (
	"reflect"
	"testing"
)

func TestIntResponse(t *testing.T) {
	tests := []struct {
		expected []byte
		input    int64
	}{
		{[]byte(":42\r\n"), 42},
		{[]byte(":-123\r\n"), -123},
		{[]byte(":0\r\n"), 0},
	}

	for _, test := range tests {
		result := IntResponse(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("IntResponse(%d) => got %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestSimpleStringResponse(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"Hello", []byte("+Hello\r\n")},
		{"World", []byte("+World\r\n")},
		{"", []byte("+\r\n")},
	}

	for _, test := range tests {
		result := SimpleStringResponse(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SimpleStringResponse(%s) => got %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestSimpleErrorResponse(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"Error message", []byte("-Error message\r\n")},
		{"Another error", []byte("-Another error\r\n")},
		{"", []byte("-\r\n")},
	}

	for _, test := range tests {
		result := SimpleErrorResponse(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SimpleErrorResponse(%s) => got %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestBulkStringResponse(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"Data", []byte("$4\r\nData\r\n")},
		{"", []byte("$0\r\n\r\n")},
	}

	for _, test := range tests {
		result := BulkStringResponse(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("BulkStringResponse(%s) => got %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestBooleanResponse(t *testing.T) {
	tests := []struct {
		expected []byte
		input    bool
	}{
		{[]byte("#1\r\n"), true},
		{[]byte("#0\r\n"), false},
	}

	for _, test := range tests {
		result := BooleanResponse(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("BooleanResponse(%t) => got %v, want %v", test.input, result, test.expected)
		}
	}
}
