package models_test

import (
	"testing"
	"time"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestPasswordReset_IsExpired(t *testing.T) {
	t.Run("should return false when token is not expired", func(t *testing.T) {
		pr := testmodels.NewPasswordReset("user@example.com", "valid-token", 
		time.Now().Add(1*time.Hour))
		assert.False(t, pr.IsExpired())
	})

	t.Run("should return true when token is expired", func(t *testing.T) {
		pr := testmodels.NewPasswordReset("user@example.com", "expired-token",
		 time.Now().Add(-1*time.Hour))
		assert.True(t, pr.IsExpired())
	})
}
