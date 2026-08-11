package utils_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"Keeps Alphanumerics", "sample1", "sample1"},
		{"Replaces Spaces", "sample one", "sample_one"},
		{"Replaces Slashes", "a/b\\c", "a_b_c"},
		{"Replaces Special Chars", "a:b*c?d", "a_b_c_d"},
		{"Trims Outer Spaces", "  sample  ", "sample"},
		{"Empty String", "", ""},
		{"All Special Chars", " /\\:*?\"<>| ", ""},
		{"Consecutive Spaces", "a  b", "a__b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utils.SanitizeFilename(tt.in))
		})
	}
}
