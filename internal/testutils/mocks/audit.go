package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockAuditRepository struct {
	GetAuditLogsFunc func(ctx context.Context, auditFilter models.AuditFilter) (
		[]models.Audit, error)
	CreateAuditFunc func(ctx context.Context, audit *models.Audit) error
}

func (r *MockAuditRepository) GetAuditLogs(ctx context.Context,
	auditFilter models.AuditFilter) ([]models.Audit, error) {
	if r.GetAuditLogsFunc != nil {
		return r.GetAuditLogsFunc(ctx, auditFilter)
	}

	return nil, nil
}

func (r *MockAuditRepository) CreateAudit(ctx context.Context,
	audit *models.Audit) error {
	if r.CreateAuditFunc != nil {
		return r.CreateAuditFunc(ctx, audit)
	}

	return nil
}

type MockAuditService struct {
	FindAllFunc func(ctx context.Context, auditFilter models.AuditFilter) (
		[]models.AuditResponse, error)
	FindAuditSelectOptionsFunc func(ctx context.Context) (
		*models.AuditSelectOptionsResponse, error)
	CreateFunc func(ctx context.Context, input *models.AuditInput) error
}

func (s *MockAuditService) FindAuditSelectOptions(ctx context.Context) (
	*models.AuditSelectOptionsResponse, error) {
	if s.FindAuditSelectOptionsFunc != nil {
		return s.FindAuditSelectOptionsFunc(ctx)
	}

	return nil, nil
}

func (s *MockAuditService) FindAll(ctx context.Context,
	auditFilter models.AuditFilter) ([]models.AuditResponse, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx, auditFilter)
	}

	return nil, nil
}

func (s *MockAuditService) Create(ctx context.Context,
	input *models.AuditInput) error {
	if s.CreateFunc != nil {
		return s.CreateFunc(ctx, input)
	}

	return nil
}
