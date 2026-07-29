package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockEmailUpdateRepository struct {
	CreateRequestFunc          func(ctx context.Context, req *models.EmailUpdateRequest) error
	GetByTokenFunc             func(ctx context.Context, token string) (*models.EmailUpdateRequest, error)
	DeleteRequestsByUserIDFunc func(ctx context.Context, userID uuid.UUID) error
}

func (r *MockEmailUpdateRepository) CreateRequest(ctx context.Context,
	req *models.EmailUpdateRequest) error {
	if r.CreateRequestFunc != nil {
		return r.CreateRequestFunc(ctx, req)
	}
	return nil
}

func (r *MockEmailUpdateRepository) GetByToken(ctx context.Context,
	token string) (*models.EmailUpdateRequest, error) {
	if r.GetByTokenFunc != nil {
		return r.GetByTokenFunc(ctx, token)
	}
	return nil, nil
}

func (r *MockEmailUpdateRepository) DeleteRequestsByUserID(ctx context.Context,
	userID uuid.UUID) error {
	if r.DeleteRequestsByUserIDFunc != nil {
		return r.DeleteRequestsByUserIDFunc(ctx, userID)
	}
	return nil
}
