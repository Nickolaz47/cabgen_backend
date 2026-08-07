package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockMetricsService struct {
	GetMetricsFunc func(ctx context.Context) (*models.AdminMetricsResponse,
		error)
}

func (s *MockMetricsService) GetMetrics(ctx context.Context) (
	*models.AdminMetricsResponse, error) {
	if s.GetMetricsFunc != nil {
		return s.GetMetricsFunc(ctx)
	}

	return nil, nil
}
