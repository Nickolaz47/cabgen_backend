package models_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestAuditToResponse(t *testing.T) {
	t.Run("should return full response when user is present", func(t *testing.T) {
		audit := testmodels.NewAudit(models.AuditEventLogin, "api",
			`{"email":"user@example.com"}`, 200)

		expected := models.AuditResponse{
			ID:        audit.ID,
			Event:     audit.Event,
			Source:    audit.Source,
			Status:    audit.Status,
			Metadata:  audit.Metadata,
			CreatedAt: audit.CreatedAt.Format(time.RFC3339),
			Username:  audit.User.Username,
		}
		result := audit.ToResponse()

		assert.Equal(t, expected, result)
	})

	t.Run("should return zero username when user is nil", func(t *testing.T) {
		audit := testmodels.NewAudit(models.AuditEventLoginFailed, "api",
			"{}", 401)
		audit.User = nil

		result := audit.ToResponse()

		assert.Empty(t, result.Username)
	})
}
