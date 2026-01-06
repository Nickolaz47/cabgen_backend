package models

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Origin struct {
	ID       string            `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Names    map[string]string `gorm:"json;not null" json:"names"`
	IsActive bool              `gorm:"not null" json:"is_active"`
}

func NewOrigin(ID string, names map[string]string, isActive bool) models.Origin {
	return models.Origin{
		ID:       uuid.MustParse(ID),
		Names:    names,
		IsActive: isActive,
	}
}

type MockOriginService struct {
	FindAllFunc       func(ctx context.Context) ([]models.Origin, error)
	FindAllActiveFunc func(ctx context.Context, lang string) ([]models.OriginFormResponse, error)
	FindByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.Origin, error)
	FindByNameFunc    func(ctx context.Context, name, lang string) ([]models.Origin, error)
	CreateFunc        func(ctx context.Context, origin *models.Origin) error
	UpdateFunc        func(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error)
	DeleteFunc        func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockOriginService) FindAll(ctx context.Context) ([]models.Origin, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockOriginService) FindAllActive(ctx context.Context, lang string) ([]models.OriginFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx, lang)
	}
	return nil, nil
}

func (m *MockOriginService) FindByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockOriginService) FindByName(ctx context.Context, name, lang string) ([]models.Origin, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(ctx, name, lang)
	}
	return nil, nil
}

func (m *MockOriginService) Create(ctx context.Context, origin *models.Origin) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, origin)
	}
	return nil
}

func (m *MockOriginService) Update(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (m *MockOriginService) Delete(ctx context.Context, ID uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}
