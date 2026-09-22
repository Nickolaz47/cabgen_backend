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

func TestTicketToResponse(t *testing.T) {
	admin := testmodels.NewAdminLoginUser()
	ticket := models.Ticket{
		ID:          uuid.New(),
		Name:        "Jão",
		Email:       "jão@mail.com",
		Institution: "Fiocruz",
		Subject:     "Wrong password",
		Message:     "Cannot access my account.",
		Status:      models.TicketStatusOpen,
		CreatedAt:   time.Date(2025, 12, 31, 1, 30, 00, 00, time.UTC),
		AdminID:     &admin.ID,
		Admin:       &admin,
	}

	expected := models.TicketResponse{
		ID:          ticket.ID,
		Name:        ticket.Name,
		Email:       ticket.Email,
		Institution: ticket.Institution,
		Subject:     ticket.Subject,
		Message:     ticket.Message,
		Status:      ticket.Status,
		CreatedAt:   ticket.CreatedAt.Format(time.RFC3339),
		AdminID:     &admin.ID,
		Admin:       admin.Username,
	}
	result := ticket.ToResponse()

	assert.Equal(t, expected, result)
}

func TestTicketFilterBinding(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected models.TicketFilter
		wantErr  bool
	}{
		{
			name:  "Success - Full filter",
			query: "status=OPEN&admin=123e4567-e89b-12d3-a456-426614174000",
			expected: models.TicketFilter{
				Status: "OPEN",
				AdminID: func() *uuid.UUID {
					id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
					return &id
				}(),
			},
		},
		{
			name:     "Success - Empty query",
			query:    "",
			expected: models.TicketFilter{},
		},
		{
			name:    "Error - Invalid admin ID",
			query:   "admin=abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := testutils.BindFilter[models.TicketFilter](tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tt.expected, filter)
			}
		})
	}
}
