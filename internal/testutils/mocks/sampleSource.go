package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockSampleSourceRepository struct {
	GetSampleSourcesFunc              func(ctx context.Context) ([]models.SampleSource, error)
	GetActiveSampleSourcesFunc        func(ctx context.Context) ([]models.SampleSource, error)
	GetSampleSourceByIDFunc           func(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error)
	GetSampleSourcesByNameOrGroupFunc func(ctx context.Context, input, language string) ([]models.SampleSource, error)
	GetSampleSourceDuplicateFunc      func(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error)
	CreateSampleSourceFunc            func(ctx context.Context, sampleSource *models.SampleSource) error
	UpdateSampleSourceFunc            func(ctx context.Context, sampleSource *models.SampleSource) error
	DeleteSampleSourceFunc            func(ctx context.Context, sampleSource *models.SampleSource) error
}

func (r *MockSampleSourceRepository) GetSampleSources(ctx context.Context) ([]models.SampleSource, error) {
	if r.GetSampleSourcesFunc != nil {
		return r.GetSampleSourcesFunc(ctx)
	}

	return nil, nil
}

func (r *MockSampleSourceRepository) GetActiveSampleSources(ctx context.Context) ([]models.SampleSource, error) {
	if r.GetActiveSampleSourcesFunc != nil {
		return r.GetActiveSampleSourcesFunc(ctx)
	}

	return nil, nil
}

func (r *MockSampleSourceRepository) GetSampleSourceByID(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
	if r.GetSampleSourceByIDFunc != nil {
		return r.GetSampleSourceByIDFunc(ctx, ID)
	}

	return nil, nil
}

func (r *MockSampleSourceRepository) GetSampleSourcesByNameOrGroup(ctx context.Context, input, language string) ([]models.SampleSource, error) {
	if r.GetSampleSourcesByNameOrGroupFunc != nil {
		return r.GetSampleSourcesByNameOrGroupFunc(ctx, input, language)
	}

	return nil, nil
}

func (r *MockSampleSourceRepository) GetSampleSourceDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
	if r.GetSampleSourceDuplicateFunc != nil {
		return r.GetSampleSourceDuplicateFunc(ctx, names, ID)
	}

	return nil, nil
}

func (r *MockSampleSourceRepository) CreateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	if r.CreateSampleSourceFunc != nil {
		return r.CreateSampleSourceFunc(ctx, sampleSource)
	}

	return nil
}

func (r *MockSampleSourceRepository) UpdateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	if r.UpdateSampleSourceFunc != nil {
		return r.UpdateSampleSourceFunc(ctx, sampleSource)
	}

	return nil
}

func (r *MockSampleSourceRepository) DeleteSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	if r.DeleteSampleSourceFunc != nil {
		return r.DeleteSampleSourceFunc(ctx, sampleSource)
	}

	return nil
}

type MockSampleSourceService struct {
	FindAllFunc           func(ctx context.Context, language string) ([]models.SampleSourceAdminTableResponse, error)
	FindAllActiveFunc     func(ctx context.Context, language string) ([]models.SampleSourceFormResponse, error)
	FindByIDFunc          func(ctx context.Context, ID uuid.UUID) (*models.SampleSourceAdminDetailResponse, error)
	FindByNameOrGroupFunc func(ctx context.Context, input, language string) ([]models.SampleSourceAdminTableResponse, error)
	CreateFunc            func(ctx context.Context, input models.SampleSourceCreateInput) (*models.SampleSourceAdminDetailResponse, error)
	UpdateFunc            func(ctx context.Context, ID uuid.UUID, input models.SampleSourceUpdateInput) (*models.SampleSourceAdminDetailResponse, error)
	DeleteFunc            func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockSampleSourceService) FindAll(
	ctx context.Context,
	language string,
) ([]models.SampleSourceAdminTableResponse, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, language)
	}
	return nil, nil
}

func (m *MockSampleSourceService) FindAllActive(
	ctx context.Context,
	language string,
) ([]models.SampleSourceFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx, language)
	}
	return nil, nil
}

func (m *MockSampleSourceService) FindByID(
	ctx context.Context,
	ID uuid.UUID,
) (*models.SampleSourceAdminDetailResponse, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockSampleSourceService) FindByNameOrGroup(
	ctx context.Context,
	input, language string,
) ([]models.SampleSourceAdminTableResponse, error) {
	if m.FindByNameOrGroupFunc != nil {
		return m.FindByNameOrGroupFunc(ctx, input, language)
	}
	return nil, nil
}

func (m *MockSampleSourceService) Create(
	ctx context.Context,
	input models.SampleSourceCreateInput,
) (*models.SampleSourceAdminDetailResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockSampleSourceService) Update(
	ctx context.Context,
	ID uuid.UUID,
	input models.SampleSourceUpdateInput,
) (*models.SampleSourceAdminDetailResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (m *MockSampleSourceService) Delete(
	ctx context.Context,
	ID uuid.UUID,
) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}
