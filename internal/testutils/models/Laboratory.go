package models

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Laboratory struct {
	ID           string `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Name         string `gorm:"not null" json:"name"`
	Abbreviation string `gorm:"not null" json:"abbreviation"`
	IsActive     bool   `gorm:"not null" json:"is_active"`
}

func NewLaboratory(ID, name, abbreviation string, isActive bool) models.Laboratory {
	return models.Laboratory{
		ID:           uuid.MustParse(ID),
		Name:         name,
		Abbreviation: abbreviation,
		IsActive:     isActive,
	}
}

type MockLaboratoryService struct {
	FindAllFunc                  func(ctx context.Context) ([]models.Laboratory, error)
	FindAllActiveFunc            func(ctx context.Context) ([]models.LaboratoryFormResponse, error)
	FindByIDFunc                 func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error)
	FindByNameOrAbbreviationFunc func(ctx context.Context, input string) ([]models.Laboratory, error)
	CreateFunc                   func(ctx context.Context, lab *models.Laboratory) error
	UpdateFunc                   func(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error)
	DeleteFunc                   func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockLaboratoryService) FindAll(ctx context.Context) ([]models.Laboratory, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindAllActive(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error) {
	if m.FindByNameOrAbbreviationFunc != nil {
		return m.FindByNameOrAbbreviationFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockLaboratoryService) Create(ctx context.Context, lab *models.Laboratory) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, lab)
	}
	return nil
}

func (m *MockLaboratoryService) Update(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (m *MockLaboratoryService) Delete(ctx context.Context, ID uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}
