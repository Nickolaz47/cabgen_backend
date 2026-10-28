package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestTicketFindAll(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)

	t.Run("Success", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketsFunc: func(ctx context.Context, filter models.TicketFilter) (
				[]models.Ticket, error) {
				return []models.Ticket{ticket}, nil
			},
		}

		service := services.NewTicketService(ticketRepo, nil, nil)
		result, err := service.FindAll(ctx, models.TicketFilter{})

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, ticket.ToResponse(), result[0])
	})

	t.Run("Error", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketsFunc: func(ctx context.Context, filter models.TicketFilter) (
				[]models.Ticket, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.FindAll(ctx, models.TicketFilter{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestTicketFindByID(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)

	t.Run("Success", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &ticket, nil
			},
		}

		service := services.NewTicketService(ticketRepo, nil, nil)
		result, err := service.FindByID(ctx, ticket.ID)

		assert.NoError(t, err)
		assert.Equal(t, ticket.ToResponse(), *result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.FindByID(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())

		var loggedTicketID string
		for _, entry := range logs.All() {
			for _, field := range entry.Context {
				if field.Key == "ticket_id" {
					loggedTicketID = field.String
				}
			}
		}
		assert.Equal(t, ticket.ID.String(), loggedTicketID)
	})

	t.Run("Error - Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.FindByID(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestTicketCreate(t *testing.T) {
	ctx := context.Background()
	input := models.CreateTicketInput{
		Name:        "Jão",
		Email:       "jão@mail.com",
		Institution: "Fiocruz",
		Subject:     "Wrong password",
		Message:     "Cannot access my account.",
	}

	t.Run("Success", func(t *testing.T) {
		createOrder, enqueueOrder := 0, 0
		ticketRepo := &mocks.MockTicketRepository{
			CreateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				createOrder = enqueueOrder + 1
				return nil
			},
		}

		asynqClient := &mocks.MockTaskQueue{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				enqueueOrder = createOrder + 1
				return &asynq.TaskInfo{ID: "task-id", Queue: "email"}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)

		service := services.NewTicketService(ticketRepo, asynqClient,
			mockLogger)
		result, err := service.Create(ctx, input, "en")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, models.TicketStatusOpen, result.Status)
		assert.NotZero(t, result.ID,
			"ticket ID must be pre-generated, not from the DB write")
		assert.Less(t, enqueueOrder, createOrder,
			"task must be enqueued before the ticket is created")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			CreateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				return gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		asynqClient := &mocks.MockTaskQueue{}

		service := services.NewTicketService(ticketRepo, asynqClient,
			mockLogger)
		result, err := service.Create(ctx, input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
		assert.True(t, asynqClient.DeleteCalled)
		assert.Equal(t, tasks.QueueEmail, asynqClient.DeleteQueue)
		assert.Equal(t, "mock-task-id", asynqClient.DeleteWith)
	})

	t.Run("Error - Enqueue Fail", func(t *testing.T) {
		created := false
		ticketRepo := &mocks.MockTicketRepository{
			CreateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				created = true
				return nil
			},
		}
		asynqClient := &mocks.MockTaskQueue{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				return nil, errors.New("redis connection refused")
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, asynqClient,
			mockLogger)
		result, err := service.Create(ctx, input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.False(t, created, "ticket must not be created when enqueue fails")
		assert.False(t, asynqClient.DeleteCalled)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestTicketAssign(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		nil,
	)

	t.Run("Success", func(t *testing.T) {
		dbTicket := ticket

		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				if dbTicket.AdminID != nil {
					dbTicket.Admin = &admin
				}
				ticketToReturn := dbTicket
				return &ticketToReturn, nil
			},
			UpdateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				dbTicket = *ticket
				return nil
			},
		}

		service := services.NewTicketService(ticketRepo, nil, nil)
		result, err := service.Assign(ctx, ticket.ID, admin.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.TicketStatusInProgress, result.Status)
		assert.Equal(t, &admin.ID, result.AdminID)
		assert.Equal(t, admin.Username, result.Admin)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.Assign(ctx, ticket.ID, admin.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - GetByID Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.Assign(ctx, ticket.ID, admin.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Invalid Ticket Status", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				inProgressTicket := ticket
				inProgressTicket.Status = models.TicketStatusInProgress
				return &inProgressTicket, nil
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.Assign(ctx, ticket.ID, admin.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrTicketIsNotOpen)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &ticket, nil
			},
			UpdateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				return gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.Assign(ctx, ticket.ID, admin.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestTicketResolve(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)

	t.Run("Success", func(t *testing.T) {
		updateOrder, enqueueOrder := 0, 0
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				ticketToResolve := ticket
				return &ticketToResolve, nil
			},
			UpdateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				updateOrder = enqueueOrder + 1
				return nil
			},
		}
		asynqClient := &mocks.MockTaskQueue{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				enqueueOrder = updateOrder + 1
				return &asynq.TaskInfo{ID: "task-id", Queue: "email"}, nil
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)

		service := services.NewTicketService(ticketRepo, asynqClient,
			mockLogger)
		result, err := service.Resolve(ctx, ticket.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.TicketStatusResolved, result.Status)
		assert.Less(t, enqueueOrder, updateOrder,
			"task must be enqueued before the ticket is updated")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Enqueue Fail", func(t *testing.T) {
		updated := false
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				ticketToResolve := ticket
				return &ticketToResolve, nil
			},
			UpdateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				updated = true
				return nil
			},
		}
		failingEnqueuer := &mocks.MockTaskQueue{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				return nil, errors.New("redis timeout")
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, failingEnqueuer,
			mockLogger)
		result, err := service.Resolve(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.False(t, updated, "ticket must not be resolved when enqueue fails")
		assert.False(t, failingEnqueuer.DeleteCalled)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.Resolve(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - GetByID Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.Resolve(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Already Resolved", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				resolvedTicket := ticket
				resolvedTicket.Status = models.TicketStatusResolved
				return &resolvedTicket, nil
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		result, err := service.Resolve(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrTicketAlreadyResolvedStatus)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &ticket, nil
			},
			UpdateTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				return gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		taskQueue := &mocks.MockTaskQueue{}
		service := services.NewTicketService(ticketRepo, taskQueue, mockLogger)
		result, err := service.Resolve(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
		assert.True(t, taskQueue.DeleteCalled)
		assert.Equal(t, tasks.QueueEmail, taskQueue.DeleteQueue)
		assert.Equal(t, "mock-task-id", taskQueue.DeleteWith)
	})
}

func TestTicketDelete(t *testing.T) {
	ctx := context.Background()
	admin := testmodels.NewAdminLoginUser()

	ticket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)

	t.Run("Success", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &ticket, nil
			},
		}

		service := services.NewTicketService(ticketRepo, nil, nil)
		err := service.Delete(ctx, ticket.ID)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		err := service.Delete(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - GetByID Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		err := service.Delete(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Ticket In Progress", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				inProgressTicket := ticket
				inProgressTicket.Status = models.TicketStatusInProgress
				return &inProgressTicket, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		err := service.Delete(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrDeleteActiveTicket)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &ticket, nil
			},
			DeleteTicketFunc: func(ctx context.Context,
				ticket *models.Ticket) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewTicketService(ticketRepo, nil, mockLogger)
		err := service.Delete(ctx, ticket.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
