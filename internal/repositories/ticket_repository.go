package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	GetTickets(ctx context.Context, status string) ([]models.Ticket, error)
	GetTicketByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
	CreateTicket(ctx context.Context, ticket *models.Ticket) error
	UpdateTicket(ctx context.Context, ticket *models.Ticket) error
	DeleteTicket(ctx context.Context, ticket *models.Ticket) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepo(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) GetTickets(ctx context.Context, status string) (
	[]models.Ticket, error) {
	var tickets []models.Ticket
	query := r.db.WithContext(ctx)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Preload("Admin").Order("created_at ASC").Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetTicketByID(ctx context.Context,
	id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.db.WithContext(ctx).Preload("Admin").First(
		&ticket, "id = ?", id).Error
	return &ticket, err
}

func (r *ticketRepository) CreateTicket(ctx context.Context,
	ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) UpdateTicket(ctx context.Context,
	ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ticketRepository) DeleteTicket(ctx context.Context,
	ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Delete(ticket).Error
}
