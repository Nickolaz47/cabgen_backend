package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type MockPasswordResetRepository struct {
	CreateTokenFunc         func(ctx context.Context, reset *models.PasswordReset) error
	GetByTokenFunc          func(ctx context.Context, token string) (*models.PasswordReset, error)
	DeleteTokensByEmailFunc func(ctx context.Context, email string) error
}

func (r *MockPasswordResetRepository) CreateToken(ctx context.Context, reset *models.PasswordReset) error {
	if r.CreateTokenFunc != nil {
		return r.CreateTokenFunc(ctx, reset)
	}
	return nil
}

func (r *MockPasswordResetRepository) GetByToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	if r.GetByTokenFunc != nil {
		return r.GetByTokenFunc(ctx, token)
	}
	return nil, nil
}

func (r *MockPasswordResetRepository) DeleteTokensByEmail(ctx context.Context, email string) error {
	if r.DeleteTokensByEmailFunc != nil {
		return r.DeleteTokensByEmailFunc(ctx, email)
	}
	return nil
}
