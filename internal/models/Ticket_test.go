package models_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTicketToResponse(t *testing.T) {
	ticket := models.Ticket{
		ID:          uuid.New(),
		Name:        "Jão",
		Email:       "jão@mail.com",
		Institution: "Fiocruz",
		Subject:     "Wrong password",
		Message:     "Cannot access my account.",
		Status:      models.TicketStatusOpen,
		CreatedAt:   time.Date(2025, 12, 31, 1, 30, 00, 00, time.UTC),
	}

	expected := models.TicketResponse{
		ID:          ticket.ID,
		Name:        ticket.Name,
		Email:       ticket.Email,
		Institution: ticket.Institution,
		Subject:     ticket.Subject,
		Message:     ticket.Message,
		Status:      ticket.Status,
		CreatedAt:   ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	result := ticket.ToResponse()

	assert.Equal(t, expected, result)
}
