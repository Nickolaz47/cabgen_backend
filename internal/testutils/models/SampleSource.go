package models

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type SampleSource struct {
	ID       string            `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Names    map[string]string `gorm:"json;not null" json:"names"`
	Groups   map[string]string `gorm:"json;not null" json:"groups"`
	IsActive bool              `gorm:"not null" json:"is_active"`
}

func NewSampleSource(ID string, names, groups map[string]string, isActive bool) models.SampleSource {
	return models.SampleSource{
		ID:       uuid.MustParse(ID),
		Names:    names,
		Groups:   groups,
		IsActive: isActive,
	}
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
