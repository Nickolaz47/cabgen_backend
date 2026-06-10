package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestSendAdminAlertEmail(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	newUser := models.User{Username: "newguy", Email: "newguy@mail.com"}
	adminUser := testmodels.NewAdminLoginUser()

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.User, error) {
				return &newUser, nil
			},
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{adminUser}, nil
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, _ := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewEmailService(userRepo, nil, nil, sender, mockLogger)
		err := svc.SendAdminAlertEmail(ctx, userID)

		assert.NoError(t, err)
	})

	t.Run("Error - Fetch New User", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, nil, sender, mockLogger)
		err := svc.SendAdminAlertEmail(ctx, userID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to fetch new user:")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Get Admins", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.User, error) {
				return &newUser, nil
			},
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, nil, sender, mockLogger)
		err := svc.SendAdminAlertEmail(ctx, userID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to get admins:")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Send Email Soft Fail", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.User, error) {
				return &newUser, nil
			},
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{adminUser}, nil
			},
		}
		sender := &mocks.MockEmailSender{
			ShouldFail: true,
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, nil, sender, mockLogger)
		err := svc.SendAdminAlertEmail(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSendWelcomeEmail(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	user := models.User{Name: "John Doe", Email: "john@mail.com"}

	t.Run("Success", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.User, error) {
				return &user, nil
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, _ := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewEmailService(userRepo, nil, nil, sender,
			mockLogger)
		err := svc.SendWelcomeEmail(ctx, userID)

		assert.NoError(t, err)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, nil, sender,
			mockLogger)
		err := svc.SendWelcomeEmail(ctx, userID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to fetch user:")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Send Email Failure", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.User, error) {
				return &user, nil
			},
		}
		sender := &mocks.MockEmailSender{
			ShouldFail: true,
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, nil, sender,
			mockLogger)
		err := svc.SendWelcomeEmail(ctx, userID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to send welcome email to")
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSendAnalysisDoneEmail(t *testing.T) {
	ctx := context.Background()
	analysisID := uuid.New()
	mockAnalysis := testmodels.CreateMockAnalysis()

	t.Run("Success - Completed", func(t *testing.T) {
		mockAnalysis.Status = models.AnalysisStatusDone
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Analysis, error) {
				return &mockAnalysis, nil
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, _ := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewEmailService(nil, analysisRepo, nil, sender,
			mockLogger)
		err := svc.SendAnalysisDoneEmail(ctx, analysisID)

		assert.NoError(t, err)
	})

	t.Run("Success - Failed Status", func(t *testing.T) {
		mockAnalysis.Status = models.AnalysisStatusFailed
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Analysis, error) {
				return &mockAnalysis, nil
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, _ := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewEmailService(nil, analysisRepo, nil, sender,
			mockLogger)
		err := svc.SendAnalysisDoneEmail(ctx, analysisID)

		assert.NoError(t, err)
	})

	t.Run("Error - Analysis Not Found", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Analysis, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(nil, analysisRepo, nil, sender,
			mockLogger)
		err := svc.SendAnalysisDoneEmail(ctx, analysisID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to fetch analysis:")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Send Mail Failure", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Analysis, error) {
				return &mockAnalysis, nil
			},
		}
		sender := &mocks.MockEmailSender{
			ShouldFail: true,
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(nil, analysisRepo, nil, sender,
			mockLogger)
		err := svc.SendAnalysisDoneEmail(ctx, analysisID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to send analysis email to")
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSendAdminTicketEmail(t *testing.T) {
	ctx := context.Background()
	ticketID := uuid.New()
	mockTicket := models.Ticket{
		Name:    "John Doe",
		Email:   "john@mail.com",
		Subject: "Test Subject",
		Message: "Test message body",
	}
	adminUser := testmodels.NewAdminLoginUser()

	t.Run("Success", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &mockTicket, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{adminUser}, nil
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, _ := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewEmailService(userRepo, nil, ticketRepo,
			sender, mockLogger)
		err := svc.SendAdminTicketEmail(ctx, ticketID)

		assert.NoError(t, err)
	})

	t.Run("Error - Fetch Ticket", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		userRepo := &mocks.MockUserRepository{}
		sender := &mocks.MockEmailSender{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, ticketRepo, sender,
			mockLogger)
		err := svc.SendAdminTicketEmail(ctx, ticketID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to fetch ticket:")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Get Admins", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &mockTicket, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, ticketRepo, sender,
			mockLogger)
		err := svc.SendAdminTicketEmail(ctx, ticketID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to get admins:")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Send Email Soft Fail", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &mockTicket, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{adminUser}, nil
			},
		}
		sender := &mocks.MockEmailSender{
			ShouldFail: true,
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(userRepo, nil, ticketRepo, sender,
			mockLogger)
		err := svc.SendAdminTicketEmail(ctx, ticketID)

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSendFinishedTicketEmail(t *testing.T) {
	ctx := context.Background()
	ticketID := uuid.New()
	mockTicket := models.Ticket{
		Name:    "Jane Doe",
		Email:   "jane@mail.com",
		Subject: "Login Issue",
		Message: "I cannot log in to the system",
	}

	t.Run("Success", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &mockTicket, nil
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, _ := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewEmailService(nil, nil, ticketRepo, sender,
			mockLogger)
		err := svc.SendFinishedTicketEmail(ctx, ticketID)

		assert.NoError(t, err)
	})

	t.Run("Error - Fetch Ticket", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		sender := &mocks.MockEmailSender{}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(nil, nil, ticketRepo, sender,
			mockLogger)
		err := svc.SendFinishedTicketEmail(ctx, ticketID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to fetch ticket:")
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Send Email Failure", func(t *testing.T) {
		ticketRepo := &mocks.MockTicketRepository{
			GetTicketByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.Ticket, error) {
				return &mockTicket, nil
			},
		}
		sender := &mocks.MockEmailSender{
			ShouldFail: true,
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewEmailService(nil, nil, ticketRepo, sender,
			mockLogger)
		err := svc.SendFinishedTicketEmail(ctx, ticketID)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to send ticket email to")
		assert.Equal(t, 1, logs.Len())
	})
}
