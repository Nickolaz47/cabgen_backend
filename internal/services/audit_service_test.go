package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestNewAuditService(t *testing.T) {
	auditRepo := &mocks.MockAuditRepository{}
	auditService := services.NewAuditService(auditRepo, nil)

	assert.NotEmpty(t, auditService)
}

func TestAuditFindAll(t *testing.T) {

	auditWithUser := testmodels.NewAudit(models.AuditEventLogin, "10.0.0.1", "{}", 200)
	auditWithoutUser := testmodels.NewAudit(
		models.AuditEventLoginFailed, "10.0.0.1", "{}", 401)
	auditWithoutUser.UserID = nil
	auditWithoutUser.User = nil

	t.Run("Success - Without Filter", func(t *testing.T) {
		auditRepo := &mocks.MockAuditRepository{
			GetAuditLogsFunc: func(ctx context.Context, _ models.AuditFilter) (
				[]models.Audit, error) {
				return []models.Audit{auditWithUser, auditWithoutUser}, nil
			},
		}

		service := services.NewAuditService(auditRepo, nil)
		auditLogs, err := service.FindAll(context.Background(),
			models.AuditFilter{})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 2)
		assert.Equal(t, auditWithUser.ToResponse(), auditLogs[0])
		assert.Equal(t, auditWithoutUser.ToResponse(), auditLogs[1])
	})

	t.Run("Success - With Filter", func(t *testing.T) {
		var receivedFilter models.AuditFilter
		auditRepo := &mocks.MockAuditRepository{
			GetAuditLogsFunc: func(ctx context.Context,
				auditFilter models.AuditFilter) ([]models.Audit, error) {
				receivedFilter = auditFilter
				return []models.Audit{auditWithUser}, nil
			},
		}

		service := services.NewAuditService(auditRepo, nil)
		filter := models.AuditFilter{Event: models.AuditEventLogin}
		auditLogs, err := service.FindAll(context.Background(), filter)

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 1)
		assert.Equal(t, filter, receivedFilter)
	})

	t.Run("Error", func(t *testing.T) {
		auditRepo := &mocks.MockAuditRepository{
			GetAuditLogsFunc: func(ctx context.Context, _ models.AuditFilter) (
				[]models.Audit, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAuditService(auditRepo, mockLogger)
		auditLogs, err := service.FindAll(context.Background(),
			models.AuditFilter{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, auditLogs)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAuditCreate(t *testing.T) {
	ctx := context.Background()

	userID := testmodels.NewAdminLoginUser().ID

	t.Run("Success - With User and Metadata", func(t *testing.T) {
		metadata := `{"auth_identity":"nick@mail.com"}`
		var receivedAudit models.Audit
		auditRepo := &mocks.MockAuditRepository{
			CreateAuditFunc: func(ctx context.Context, audit *models.Audit) error {
				receivedAudit = *audit
				return nil
			},
		}

		service := services.NewAuditService(auditRepo, nil)
		err := service.Create(ctx, &models.AuditInput{
			Event:    models.AuditEventLoginFailed,
			Source:   "10.0.0.1",
			Status:   401,
			Metadata: &metadata,
			UserID:   &userID,
		})

		assert.NoError(t, err)
		assert.Equal(t, models.AuditEventLoginFailed, receivedAudit.Event)
		assert.Equal(t, "10.0.0.1", receivedAudit.Source)
		assert.Equal(t, 401, receivedAudit.Status)
		assert.Equal(t, metadata, receivedAudit.Metadata)
		assert.Equal(t, &userID, receivedAudit.UserID)
	})

	t.Run("Success - Without User and Metadata", func(t *testing.T) {
		var receivedAudit models.Audit
		auditRepo := &mocks.MockAuditRepository{
			CreateAuditFunc: func(ctx context.Context, audit *models.Audit) error {
				receivedAudit = *audit
				return nil
			},
		}

		service := services.NewAuditService(auditRepo, nil)
		err := service.Create(ctx, &models.AuditInput{
			Event:  models.AuditEventLogin,
			Source: "10.0.0.1",
			Status: 200,
		})

		assert.NoError(t, err)
		assert.Empty(t, receivedAudit.Metadata)
		assert.Nil(t, receivedAudit.UserID)
	})

	t.Run("Error", func(t *testing.T) {
		auditRepo := &mocks.MockAuditRepository{
			CreateAuditFunc: func(ctx context.Context, audit *models.Audit) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		service := services.NewAuditService(auditRepo, mockLogger)
		err := service.Create(ctx, &models.AuditInput{
			Event:  models.AuditEventLogin,
			Source: "10.0.0.1",
			Status: 200,
		})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
