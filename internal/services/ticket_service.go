package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TicketService interface {
	FindAll(ctx context.Context, filter models.TicketFilter) (
		[]models.TicketResponse, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*models.TicketResponse, error)
	Create(ctx context.Context, input models.CreateTicketInput, language string) (
		*models.TicketResponse, error)
	Assign(ctx context.Context, ticketID uuid.UUID, adminID uuid.UUID) (
		*models.TicketResponse, error)
	Resolve(ctx context.Context, ticketID uuid.UUID) (*models.TicketResponse,
		error)
	Delete(ctx context.Context, ticketID uuid.UUID) error
}

type ticketService struct {
	Repo        repositories.TicketRepository
	AsynqClient TaskQueue
	Logger      *zap.Logger
}

func NewTicketService(
	repo repositories.TicketRepository,
	asynqClient TaskQueue,
	logger *zap.Logger,
) TicketService {
	return &ticketService{
		Repo:        repo,
		AsynqClient: asynqClient,
		Logger:      logger,
	}
}

func (s *ticketService) FindAll(ctx context.Context, filter models.TicketFilter) (
	[]models.TicketResponse, error) {
	tickets, err := s.Repo.GetTickets(ctx, filter)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
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

func (s *ticketService) FindByID(ctx context.Context, ID uuid.UUID) (
	*models.TicketResponse, error) {
	ticket, err := s.Repo.GetTicketByID(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Warn("Service Warning", logging.ServiceLogging(ctx,
				"TicketService", "FindByID", logging.DatabaseNotFoundError, err,
				zap.String("ticket_id", ID.String()),
			)...)
			return nil, ErrNotFound
		}

		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "FindByID", logging.DatabaseError, err,
			zap.String("ticket_id", ID.String()),
		)...)
		return nil, ErrInternal
	}

	response := ticket.ToResponse()
	return &response, nil
}

func (s *ticketService) Create(ctx context.Context,
	input models.CreateTicketInput, language string) (*models.TicketResponse, error) {
	ticket := models.Ticket{
		Name:        input.Name,
		Email:       input.Email,
		Institution: input.Institution,
		Subject:     input.Subject,
		Message:     input.Message,
		Status:      models.TicketStatusOpen,
		Language:    translation.ParseLanguage(language),
	}

	// Pre-generated so the task can reference the ticket before insertion
	id, err := uuid.NewV7()
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Create", logging.GenerateIDError, err,
		)...)
		return nil, ErrInternal
	}
	ticket.ID = id

	task, err := tasks.NewAdminTicketEmailTask(ticket.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Create", logging.AsynqTaskError, err,
		)...)
		return nil, ErrInternal
	}

	info, err := s.AsynqClient.EnqueueContext(ctx, task,
		asynq.Queue(tasks.QueueEmail), asynq.ProcessIn(CommitGrace))
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Create", logging.RedisDispatchError, err,
		)...)
		return nil, ErrInternal
	}
	s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(ctx,
		"TicketService", "Create", logging.TaskEnqueuedSuccess,
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
	)...)

	if err := s.Repo.CreateTicket(ctx, &ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Create", logging.DatabaseError, err,
		)...)
		DeleteScheduledTask(ctx, s.AsynqClient, s.Logger,
			"TicketService", "Create", tasks.QueueEmail, info.ID)
		return nil, ErrInternal
	}

	response := ticket.ToResponse()
	return &response, nil
}

func (s *ticketService) Assign(ctx context.Context, ticketID,
	adminID uuid.UUID) (*models.TicketResponse, error) {
	ticket, err := s.Repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Warn("Service Warning", logging.ServiceLogging(ctx,
				"TicketService", "Assign", logging.DatabaseNotFoundError, err,
				zap.String("ticket_id", ticketID.String()),
			)...)
			return nil, ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Assign", logging.DatabaseError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return nil, ErrInternal
	}

	if ticket.Status != models.TicketStatusOpen {
		s.Logger.Warn("Business Rule Violation", logging.ServiceLogging(ctx,
			"TicketService", "Assign", logging.TicketStatusError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return nil, ErrTicketIsNotOpen
	}

	ticket.AdminID = &adminID
	ticket.Status = models.TicketStatusInProgress

	if err := s.Repo.UpdateTicket(ctx, ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Assign", logging.DatabaseError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return nil, ErrInternal
	}

	updatedTicket, err := s.Repo.GetTicketByID(ctx, ticket.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Assign", logging.DatabaseError, err,
			zap.String("ticket_id", ticketID.String()),
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
			s.Logger.Warn("Service Warning", logging.ServiceLogging(ctx,
				"TicketService", "Resolve", logging.DatabaseNotFoundError, err,
				zap.String("ticket_id", ticketID.String()),
			)...)
			return nil, ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Resolve", logging.DatabaseError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return nil, ErrInternal
	}

	if ticket.Status == models.TicketStatusResolved {
		s.Logger.Warn("Business Rule Violation", logging.ServiceLogging(ctx,
			"TicketService", "Resolve", logging.TicketStatusError,
			ErrTicketAlreadyResolvedStatus)...)
		return nil, ErrTicketAlreadyResolvedStatus
	}

	ticket.Status = models.TicketStatusResolved

	task, err := tasks.NewFinishedTicketEmailTask(ticket.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Resolve", logging.AsynqTaskError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return nil, ErrInternal
	}

	info, err := s.AsynqClient.EnqueueContext(ctx, task,
		asynq.Queue(tasks.QueueEmail), asynq.ProcessIn(CommitGrace))
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Resolve", logging.RedisDispatchError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return nil, ErrInternal
	}
	s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(ctx,
		"TicketService", "Resolve", logging.TaskEnqueuedSuccess,
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.String("ticket_id", ticketID.String()),
	)...)

	if err := s.Repo.UpdateTicket(ctx, ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Resolve", logging.DatabaseError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		DeleteScheduledTask(ctx, s.AsynqClient, s.Logger,
			"TicketService", "Resolve", tasks.QueueEmail, info.ID)
		return nil, ErrInternal
	}

	response := ticket.ToResponse()
	return &response, nil
}

func (s *ticketService) Delete(ctx context.Context, ticketID uuid.UUID) error {
	ticket, err := s.Repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Warn("Service Warning", logging.ServiceLogging(ctx,
				"TicketService", "Delete", logging.DatabaseNotFoundError, err,
				zap.String("ticket_id", ticketID.String()),
			)...)
			return ErrNotFound
		}
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Delete", logging.DatabaseError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return ErrInternal
	}

	if ticket.Status == models.TicketStatusInProgress {
		err := errors.New("cannot delete a ticket that is currently in progress")
		s.Logger.Warn("Business Rule Violation", logging.ServiceLogging(ctx,
			"TicketService", "Delete", logging.DeleteActiveTicketError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return ErrDeleteActiveTicket
	}

	if err := s.Repo.DeleteTicket(ctx, ticket); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"TicketService", "Delete", logging.DatabaseError, err,
			zap.String("ticket_id", ticketID.String()),
		)...)
		return ErrInternal
	}

	return nil
}
