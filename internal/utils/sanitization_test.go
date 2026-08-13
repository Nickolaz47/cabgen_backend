package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"trim spaces", "  hello  ", "hello"},
		{"trim tabs", "\t\thello\t\t", "hello"},
		{"replace newlines", "hello\nworld", "hello world"},
		{"replace carriage returns", "hello\rworld", "hello world"},
		{"replace tabs inside", "hello\tworld", "hello world"},
		{"replace mixed whitespace", "hello\r\n\tworld", "hello   world"},
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
		{"only newlines", "\n\n\n", ""},
		{"single word", "hello", "hello"},
		{"no special chars", "hello world", "hello world"},
		{"vertical tab", "hello\vworld", "hello world"},
		{"form feed", "hello\fworld", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeQuery(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
