package utils_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectRoot(t *testing.T) {
	result, err := utils.GetProjectRoot()

	assert.NoError(t, err, "Expected no error to get project root")
	assert.NotEmpty(t, result)
	assert.DirExists(t, result)
}
