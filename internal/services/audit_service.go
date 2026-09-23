package services

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"go.uber.org/zap"
)

type AuditService interface {
	FindAll(ctx context.Context, auditFilter models.AuditFilter) (
		[]models.AuditResponse, error)
	FindAuditSelectOptions(ctx context.Context) (
		*models.AuditSelectOptionsResponse, error)
	Create(ctx context.Context, input *models.AuditInput) error
}

type auditService struct {
	AuditRepo repositories.AuditRepository
	UserRepo  repositories.UserRepository
	Logger    *zap.Logger
}

func NewAuditService(
	auditRepo repositories.AuditRepository,
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) AuditService {
	return &auditService{
		AuditRepo: auditRepo,
		UserRepo:  userRepo,
		Logger:    logger,
	}
}

func (s *auditService) FindAll(ctx context.Context,
	auditFilter models.AuditFilter) ([]models.AuditResponse, error) {
	auditLogs, err := s.AuditRepo.GetAuditLogs(ctx, auditFilter)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"AuditService", "FindAll", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	responses := make([]models.AuditResponse, len(auditLogs))
	for i := range auditLogs {
		responses[i] = auditLogs[i].ToResponse()
	}
	return responses, nil
}

func (s *auditService) FindAuditSelectOptions(ctx context.Context) (
	*models.AuditSelectOptionsResponse, error) {
	users, err := s.UserRepo.GetUsers(ctx, models.AdminUserFilter{})
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"AuditService", "FindAuditSelectOptions", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	userOptions := make([]models.SelectOption, len(users))
	for i, user := range users {
		userOptions[i] = models.SelectOption{
			Label: user.Username,
			Value: user.ID.String(),
		}
	}

	return &models.AuditSelectOptionsResponse{
		Events: models.AuditEventSelectOptions(),
		Users:  userOptions,
	}, nil
}

func (s *auditService) Create(ctx context.Context,
	input *models.AuditInput) error {
	audit := models.Audit{
		Event:    input.Event,
		Source:   input.Source,
		Status:   input.Status,
		Metadata: metadataString(input.Metadata),
		UserID:   input.UserID,
	}

	if err := s.AuditRepo.CreateAudit(ctx, &audit); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"AuditService", "Create", logging.CreateAuditError, err,
		)...)
		return ErrInternal
	}

	return nil
}

func metadataString(metadata *string) string {
	if metadata == nil {
		return ""
	}
	return *metadata
}
