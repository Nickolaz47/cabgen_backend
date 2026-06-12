package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TicketService interface {
	FindAll(ctx context.Context, status string) ([]models.TicketResponse, error)
	Create(ctx context.Context, input models.CreateTicketInput) (
		*models.TicketResponse, error)
	Assign(ctx context.Context, ticketID uuid.UUID, adminID uuid.UUID) (
		*models.TicketResponse, error)
	Resolve(ctx context.Context, ticketID uuid.UUID) (*models.TicketResponse,
		error)
	Delete(ctx context.Context, ticketID uuid.UUID) error
}

type ticketService struct {
	Repo        repositories.TicketRepository
	AsynqClient TaskEnqueuer
	Logger      *zap.Logger
}

func NewTicketService(
	repo repositories.TicketRepository,
	asynqClient TaskEnqueuer,
	logger *zap.Logger,
) TicketService {
	return &ticketService{
		Repo:        repo,
		AsynqClient: asynqClient,
		Logger:      logger,
	}
}

func (s *ticketService) FindAll(ctx context.Context, status string) (
	[]models.TicketResponse, error) {
	tickets, err := s.Repo.GetTickets(ctx, status)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "FindAll", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	responses := make([]models.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = ticket.ToResponse()
	}

	return responses, nil
}

func (s *ticketService) Create(ctx context.Context,
	input models.CreateTicketInput) (*models.TicketResponse, error) {
	ticket := models.Ticket{
		Name:        input.Name,
		Email:       input.Email,
		Institution: input.Institution,
		Subject:     input.Subject,
		Message:     input.Message,
		Status:      models.TicketStatusOpen,
	}

	if err := s.Repo.CreateTicket(ctx, &ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Create", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	task, err := tasks.NewAdminTicketEmailTask(ticket.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Create", logging.AsynqTaskError, err,
		)...)
	} else {
		info, err := s.AsynqClient.EnqueueContext(ctx, task,
			asynq.Queue(tasks.QueueEmail))
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"TicketService", "Create", logging.RedisDispatchError, err,
			)...)
		} else {
			s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(
				"TicketService", "Create", logging.TaskEnqueuedSuccess,
				zap.String("task_id", info.ID),
				zap.String("queue", info.Queue),
			)...)
		}
	}

	response := ticket.ToResponse()
	return &response, nil
}

func (s *ticketService) Assign(ctx context.Context, ticketID,
	adminID uuid.UUID) (*models.TicketResponse, error) {
	ticket, err := s.Repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"TicketService", "Assign", logging.DatabaseNotFoundError, err,
			)...)
			return nil, ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Assign", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if ticket.Status != models.TicketStatusOpen {
		s.Logger.Warn("Business Rule Violation", logging.ServiceLogging(
			"TicketService", "Assign", logging.TicketStatusError, err,
		)...)
		return nil, ErrTicketIsNotOpen
	}

	ticket.AdminID = &adminID
	ticket.Status = models.TicketStatusInProgress

	if err := s.Repo.UpdateTicket(ctx, ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Assign", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	updatedTicket, err := s.Repo.GetTicketByID(ctx, ticket.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Assign", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	response := updatedTicket.ToResponse()
	return &response, nil
}

func (s *ticketService) Resolve(ctx context.Context, ticketID uuid.UUID) (
	*models.TicketResponse, error) {
	ticket, err := s.Repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"TicketService", "Resolve", logging.DatabaseNotFoundError, err,
			)...)
			return nil, ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Resolve", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if ticket.Status == models.TicketStatusResolved {
		s.Logger.Warn("Business Rule Violation", logging.ServiceLogging(
			"TicketService", "Resolve", logging.TicketStatusError,
			ErrTicketAlreadyResolvedStatus)...)
		return nil, ErrTicketAlreadyResolvedStatus
	}

	ticket.Status = models.TicketStatusResolved

	if err := s.Repo.UpdateTicket(ctx, ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Resolve", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	task, err := tasks.NewFinishedTicketEmailTask(ticket.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Resolve", logging.AsynqTaskError, err,
		)...)
	} else {
		info, err := s.AsynqClient.EnqueueContext(ctx, task,
			asynq.Queue(tasks.QueueEmail))
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"TicketService", "Resolve", logging.RedisDispatchError, err,
			)...)
		} else {
			s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(
				"TicketService", "Resolve", logging.TaskEnqueuedSuccess,
				zap.String("task_id", info.ID),
				zap.String("queue", info.Queue),
			)...)
		}
	}

	response := ticket.ToResponse()
	return &response, nil
}

func (s *ticketService) Delete(ctx context.Context, ticketID uuid.UUID) error {
	ticket, err := s.Repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"TicketService", "Delete", logging.DatabaseNotFoundError, err,
			)...)
			return ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if ticket.Status == models.TicketStatusInProgress {
		err := errors.New("cannot delete a ticket that is currently in progress")
		s.Logger.Warn("Business Rule Violation", logging.ServiceLogging(
			"TicketService", "Delete", logging.DeleteActiveTicketError, err,
		)...)
		return ErrDeleteActiveTicket
	}

	if err := s.Repo.DeleteTicket(ctx, ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"TicketService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}
