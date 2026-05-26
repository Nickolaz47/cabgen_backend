package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockSampleRepository struct {
	GetSamplesFunc func(ctx context.Context, input string,
		userID uuid.UUID) ([]models.Sample, error)
	GetSampleByIDFunc func(ctx context.Context,
		ID uuid.UUID) (*models.Sample, error)
	CreateSampleFunc func(ctx context.Context, sample *models.Sample) error
	UpdateSampleFunc func(ctx context.Context, sample *models.Sample) error
	DeleteSampleFunc func(ctx context.Context, sample *models.Sample) error
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

type MockSampleService struct {
	PrepareSampleFolderFunc func(userID, sampleID uuid.UUID) (string, error)
	FindAllFunc             func(ctx context.Context, input string,
		userID uuid.UUID, language string) ([]models.SampleResponse, error)
	FindByIDFunc func(ctx context.Context, sampleID, userID uuid.UUID,
		language string) (*models.SampleResponse, error)
	CreateFunc func(ctx context.Context, input models.SampleCreateInput,
		language string) (*models.SampleResponse, error)
	AttachFilesFunc func(ctx context.Context, sampleID, userID uuid.UUID,
		input models.SampleAttachmentInput) error
	UpdateFunc func(ctx context.Context, sampleID, userID uuid.UUID,
		input models.SampleUpdateInput,
		language string) (*models.SampleResponse, error)
	DeleteFunc func(ctx context.Context, sampleID, userID uuid.UUID) error
}

func (r *MockSampleService) PrepareSampleFolder(
	userID, sampleID uuid.UUID) (string, error) {
	if r.PrepareSampleFolderFunc != nil {
		return r.PrepareSampleFolderFunc(userID, sampleID)
	}

	return "", nil
}

func (r *MockSampleService) FindAll(ctx context.Context, input string,
	userID uuid.UUID, language string) ([]models.SampleResponse, error) {
	if r.FindAllFunc != nil {
		return r.FindAllFunc(ctx, input, userID, language)
	}

	return nil, nil
}

func (r *MockSampleService) FindByID(ctx context.Context, sampleID,
	userID uuid.UUID, language string) (*models.SampleResponse, error) {
	if r.FindByIDFunc != nil {
		return r.FindByIDFunc(ctx, sampleID, userID, language)
	}
	return nil, nil
}

func (r *MockSampleService) Create(ctx context.Context,
	input models.SampleCreateInput, language string,
) (*models.SampleResponse, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, input, language)
	}
	return nil, nil
}

func (r *MockSampleService) AttachFiles(ctx context.Context,
	sampleID, userID uuid.UUID,
	input models.SampleAttachmentInput) error {
	if r.AttachFilesFunc != nil {
		return r.AttachFilesFunc(ctx, sampleID, userID, input)
	}
	return nil
}

func (r *MockSampleService) Update(ctx context.Context, sampleID,
	userID uuid.UUID, input models.SampleUpdateInput, language string,
) (*models.SampleResponse, error) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, sampleID, userID, input, language)
	}

	return nil, nil
}

func (r *MockSampleService) Delete(ctx context.Context, sampleID,
	userID uuid.UUID) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, sampleID, userID)
	}
	return nil
}
