package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockCityService struct {
	FindAllFunc func(ctx context.Context) ([]models.SelectOption, error)
}

func (s *MockCityService) FindAll(ctx context.Context) ([]models.SelectOption,
	error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx)
	}

	return nil, nil
}
