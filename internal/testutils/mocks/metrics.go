package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockMetricsService struct {
	GetMetricsFunc func(ctx context.Context) (*models.AdminMetricsResponse,
		error)
	GetUserMetricsFunc func(ctx context.Context, userID uuid.UUID) (
		*models.UserMetricsResponse, error)
}

func (s *MockMetricsService) GetMetrics(ctx context.Context) (
	*models.AdminMetricsResponse, error) {
	if s.GetMetricsFunc != nil {
		return s.GetMetricsFunc(ctx)
	}

	return nil, nil
}

func (s *MockMetricsService) GetUserMetrics(ctx context.Context,
	userID uuid.UUID) (*models.UserMetricsResponse, error) {
	if s.GetUserMetricsFunc != nil {
		return s.GetUserMetricsFunc(ctx, userID)
	}

	return nil, nil
}
