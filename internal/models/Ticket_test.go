package models_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
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
