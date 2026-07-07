package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockSelectOptionsService struct {
	FindAllFunc func(ctx context.Context) (*models.EnumSelectsResponse, error)
}

func (s *MockSelectOptionsService) FindAll(ctx context.Context) (
	*models.EnumSelectsResponse, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx)
	}

	return nil, nil
}
