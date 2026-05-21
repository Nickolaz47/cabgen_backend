package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockSampleRepository struct {
	GetSamplesFunc func(ctx context.Context, input string,
		userID uuid.UUID) ([]models.Sample, error)
	GetSampleByIDFunc func(ctx context.Context, ID uuid.UUID) (*models.Sample, error)
	CreateSampleFunc  func(ctx context.Context, sample *models.Sample) error
	UpdateSampleFunc  func(ctx context.Context, sample *models.Sample) error
	DeleteSampleFunc  func(ctx context.Context, sample *models.Sample) error
}

func (r *MockSampleRepository) GetSamples(ctx context.Context,
	input string, userID uuid.UUID) ([]models.Sample, error) {
	if r.GetSamplesFunc != nil {
		return r.GetSamplesFunc(ctx, input, userID)
	}

	return nil, nil
}

func (r *MockSampleRepository) GetSampleByID(ctx context.Context,
	ID uuid.UUID) (*models.Sample, error) {
	if r.GetSampleByIDFunc != nil {
		return r.GetSampleByIDFunc(ctx, ID)
	}

	return nil, nil
}

func (r *MockSampleRepository) CreateSample(ctx context.Context,
	sample *models.Sample) error {
	if r.CreateSampleFunc != nil {
		return r.CreateSampleFunc(ctx, sample)
	}

	return nil
}

func (r *MockSampleRepository) UpdateSample(ctx context.Context,
	sample *models.Sample) error {
	if r.UpdateSampleFunc != nil {
		return r.UpdateSampleFunc(ctx, sample)
	}

	return nil
}

func (r *MockSampleRepository) DeleteSample(ctx context.Context,
	sample *models.Sample) error {
	if r.DeleteSampleFunc != nil {
		return r.DeleteSampleFunc(ctx, sample)
	}

	return nil
}
