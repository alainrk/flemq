package config

import (
	"os"
	"testing"
)

const TEST_PREFIX = "TEST_"

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		defaultVal any
		expected   any
		name       string
		envValue   string
	}{
		// Test cases for string
		{"default_value", "test_value", "StringSet", "test_value"},
		{"default_value", "default_value", "StringNotSet", ""},
		{"", "", "StringEmptySet", ""},

		// Test cases for bool
		{false, true, "BoolTrueSet", "true"},
		{true, false, "BoolFalseSet", "false"},
		{true, true, "BoolNotSet", ""},

		// Test cases for int
		{0, 42, "IntSet", "42"},
		{10, 10, "IntNotSet", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(TEST_PREFIX+tt.name, tt.envValue)
			defer os.Unsetenv(TEST_PREFIX + "_" + tt.name)

			result := loadEnv(TEST_PREFIX, tt.name, tt.defaultVal)

			if result != tt.expected {
				t.Errorf("For %s, expected %v, but got %v", tt.name, tt.expected, result)
			}
		})
	}
}
