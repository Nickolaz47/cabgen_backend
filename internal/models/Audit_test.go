package models_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuditToResponse(t *testing.T) {
	audit := testmodels.NewAudit(models.AuditEventLogin, "10.0.0.1",
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
}

func TestAuditToResponseWithoutUser(t *testing.T) {
	audit := testmodels.NewAudit(models.AuditEventLoginFailed, "10.0.0.1",
		"{}", 401)
	audit.User = nil

	result := audit.ToResponse()

	assert.Empty(t, result.Username)
}

func TestAuditFilterBinding(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected models.AuditFilter
		wantErr  bool
	}{
		{
			name: "Success - Full filter",
			query: "event=auth.login&source=10.0.0.1&status=404&date=2026-01-02&user=" +
				"123e4567-e89b-12d3-a456-426614174000",
			expected: models.AuditFilter{
				Event:  "auth.login",
				Source: "10.0.0.1",
				Status: 404,
				Date: func() *time.Time {
					d := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
					return &d
				}(),
				UserID: func() *uuid.UUID {
					id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
					return &id
				}(),
			},
		},
		{
			name:     "Success - Empty query",
			query:    "",
			expected: models.AuditFilter{},
		},
		{
			name:    "Error - Invalid date",
			query:   "date=02/01/2026",
			wantErr: true,
		},
		{
			name:    "Error - Invalid user ID",
			query:   "user=abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := testutils.BindFilter[models.AuditFilter](tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tt.expected, filter)
			}
		})
	}
}
