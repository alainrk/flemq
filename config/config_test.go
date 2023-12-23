package config

import (
	"os"
	"testing"
)

const TEST_PREFIX = "TEST"

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name       string
		envValue   string
		defaultVal interface{}
		expected   interface{}
	}{
		// Test cases for string
		{"StringSet", "test_value", "default_value", "test_value"},
		{"StringNotSet", "", "default_value", "default_value"},
		{"StringEmptySet", "", "", ""},

		// Test cases for bool
		{"BoolTrueSet", "true", false, true},
		{"BoolFalseSet", "false", true, false},
		{"BoolNotSet", "", true, true},

		// Test cases for int
		{"IntSet", "42", 0, 42},
		{"IntNotSet", "", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(TEST_PREFIX+"_"+tt.name, tt.envValue)
			defer os.Unsetenv(TEST_PREFIX + "_" + tt.name)

			result := loadEnv(TEST_PREFIX, tt.name, tt.defaultVal)

			if result != tt.expected {
				t.Errorf("For %s, expected %v, but got %v", tt.name, tt.expected, result)
			}
		})
	}
}
