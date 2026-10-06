package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockCityService struct {
	FindAllFunc func(ctx context.Context, language string) (
		[]models.SelectOption, error)
}

func (s *MockCityService) FindAll(ctx context.Context, language string) (
	[]models.SelectOption, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx, language)
	}

	return nil, nil
}
