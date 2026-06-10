package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockTicketRepository struct {
	GetTicketsFunc func(ctx context.Context, status string) (
		[]models.Ticket, error)
	GetTicketByIDFunc func(ctx context.Context, id uuid.UUID) (
		*models.Ticket, error)
	CreateTicketFunc func(ctx context.Context, ticket *models.Ticket) error
	UpdateTicketFunc func(ctx context.Context, ticket *models.Ticket) error
	DeleteTicketFunc func(ctx context.Context, ticket *models.Ticket) error
}

func (r *MockTicketRepository) GetTickets(ctx context.Context, status string) (
	[]models.Ticket, error) {
	if r.GetTicketsFunc != nil {
		return r.GetTicketsFunc(ctx, status)
	}
	return nil, nil
}

func (r *MockTicketRepository) GetTicketByID(ctx context.Context,
	id uuid.UUID) (*models.Ticket, error) {
	if r.GetTicketByIDFunc != nil {
		return r.GetTicketByIDFunc(ctx, id)
	}
	return nil, nil
}

func (r *MockTicketRepository) CreateTicket(ctx context.Context,
	ticket *models.Ticket) error {
	if r.CreateTicketFunc != nil {
		return r.CreateTicketFunc(ctx, ticket)
	}
	return nil
}

func (r *MockTicketRepository) UpdateTicket(ctx context.Context,
	ticket *models.Ticket) error {
	if r.UpdateTicketFunc != nil {
		return r.UpdateTicketFunc(ctx, ticket)
	}
	return nil
}

func (r *MockTicketRepository) DeleteTicket(ctx context.Context,
	ticket *models.Ticket) error {
	if r.DeleteTicketFunc != nil {
		return r.DeleteTicketFunc(ctx, ticket)
	}
	return nil
}
