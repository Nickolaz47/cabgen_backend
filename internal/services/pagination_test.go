package services_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestTotalPages(t *testing.T) {
	assert.Equal(t, 0, services.TotalPages(0))
	assert.Equal(t, 1, services.TotalPages(1))
	assert.Equal(t, 1, services.TotalPages(100))
	assert.Equal(t, 2, services.TotalPages(101))
	assert.Equal(t, 2, services.TotalPages(200))
}
