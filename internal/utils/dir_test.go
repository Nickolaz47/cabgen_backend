package utils_test

import (
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectRoot(t *testing.T) {
	result, err := utils.GetProjectRoot()

	assert.NoError(t, err, "Expected no error to get project root")

	result = filepath.Base(result)
	expected := "cabgen_backend"

	assert.Equal(t, expected, result, "Expected root directory to be %v, got %v", expected, result)
}
