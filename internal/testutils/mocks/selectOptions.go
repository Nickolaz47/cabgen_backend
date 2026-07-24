package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockSelectOptionsService struct {
	FindAllEnumSelectsFunc func(ctx context.Context) (
		*models.EnumSelectsResponse, error)
	FindAllFormSelectsFunc func(ctx context.Context, language string) (
		*models.FormSelectsResponse, error)
}

func (s *MockSelectOptionsService) FindAllEnumSelects(ctx context.Context) (
	*models.EnumSelectsResponse, error) {
	if s.FindAllEnumSelectsFunc != nil {
		return s.FindAllEnumSelectsFunc(ctx)
	}

	return nil, nil
}

func (s *MockSelectOptionsService) FindAllFormSelects(ctx context.Context,
	language string) (*models.FormSelectsResponse, error) {
	if s.FindAllFormSelectsFunc != nil {
		return s.FindAllFormSelectsFunc(ctx, language)
	}

	return nil, nil
}
