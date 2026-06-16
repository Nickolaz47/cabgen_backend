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

type MockTicketService struct {
	FindAllFunc func(ctx context.Context, status string) (
		[]models.TicketResponse, error)
	FindByIDFunc func(ctx context.Context, ID uuid.UUID) (
		*models.TicketResponse, error)
	CreateFunc func(ctx context.Context, input models.CreateTicketInput) (
		*models.TicketResponse, error)
	AssignFunc func(ctx context.Context, ticketID, adminID uuid.UUID) (
		*models.TicketResponse, error)
	ResolveFunc func(ctx context.Context, ticketID uuid.UUID) (
		*models.TicketResponse, error)
	DeleteFunc func(ctx context.Context, ticketID uuid.UUID) error
}

func (s *MockTicketService) FindAll(ctx context.Context, status string) (
	[]models.TicketResponse, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx, status)
	}
	return nil, nil
}

func (s *MockTicketService) FindByID(ctx context.Context, ID uuid.UUID) (
	*models.TicketResponse, error) {
	if s.FindByIDFunc != nil {
		return s.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (s *MockTicketService) Create(ctx context.Context,
	input models.CreateTicketInput) (*models.TicketResponse, error) {
	if s.CreateFunc != nil {
		return s.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (s *MockTicketService) Assign(ctx context.Context, ticketID,
	adminID uuid.UUID) (*models.TicketResponse, error) {
	if s.AssignFunc != nil {
		return s.AssignFunc(ctx, ticketID, adminID)
	}
	return nil, nil
}

func (s *MockTicketService) Resolve(ctx context.Context, ticketID uuid.UUID) (
	*models.TicketResponse, error) {
	if s.ResolveFunc != nil {
		return s.ResolveFunc(ctx, ticketID)
	}
	return nil, nil
}

func (s *MockTicketService) Delete(ctx context.Context,
	ticketID uuid.UUID) error {
	if s.DeleteFunc != nil {
		return s.DeleteFunc(ctx, ticketID)
	}
	return nil
}
