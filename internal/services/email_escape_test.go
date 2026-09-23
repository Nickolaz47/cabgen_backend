package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeData(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		expected map[string]any
	}{
		{
			name: "Success - Escapes script tags",
			data: map[string]any{"Name": "<script>alert(1)</script>John"},
			expected: map[string]any{
				"Name": "&lt;script&gt;alert(1)&lt;/script&gt;John"},
		},
		{
			name:     "Success - Escapes ampersand",
			data:     map[string]any{"Subject": "A & B"},
			expected: map[string]any{"Subject": "A &amp; B"},
		},
		{
			name:     "Success - Non-strings untouched",
			data:     map[string]any{"Count": 42, "Flag": true},
			expected: map[string]any{"Count": 42, "Flag": true},
		},
		{
			name:     "Success - Nil data",
			data:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, escapeData(tt.data))
		})
	}
}
