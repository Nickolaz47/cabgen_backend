package models

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Sequencer struct {
	ID       string `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Model    string `gorm:"not null" json:"model"`
	Brand    string `gorm:"not null" json:"brand"`
	IsActive bool   `gorm:"not null" json:"is_active"`
}

func NewSequencer(ID, model, brand string, isActive bool) models.Sequencer {
	return models.Sequencer{
		ID:       uuid.MustParse(ID),
		Model:    model,
		Brand:    brand,
		IsActive: isActive,
	}
}

type MockSequencerService struct {
	FindAllFunc            func(ctx context.Context) ([]models.SequencerAdminTableResponse, error)
	FindAllActiveFunc      func(ctx context.Context) ([]models.SequencerFormResponse, error)
	FindByIDFunc           func(ctx context.Context, ID uuid.UUID) (*models.SequencerAdminTableResponse, error)
	FindByBrandOrModelFunc func(ctx context.Context, input string) ([]models.SequencerAdminTableResponse, error)
	CreateFunc             func(ctx context.Context, input models.SequencerCreateInput) (*models.SequencerAdminTableResponse, error)
	UpdateFunc             func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error)
	DeleteFunc             func(ctx context.Context, ID uuid.UUID) error
}

func (s *MockSequencerService) FindAll(ctx context.Context) ([]models.SequencerAdminTableResponse, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx)
	}
	return nil, nil
}

func (s *MockSequencerService) FindAllActive(ctx context.Context) ([]models.SequencerFormResponse, error) {
	if s.FindAllActiveFunc != nil {
		return s.FindAllActiveFunc(ctx)
	}
	return nil, nil
}

func (s *MockSequencerService) FindByID(ctx context.Context, ID uuid.UUID) (*models.SequencerAdminTableResponse, error) {
	if s.FindByIDFunc != nil {
		return s.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (s *MockSequencerService) FindByBrandOrModel(ctx context.Context, input string) ([]models.SequencerAdminTableResponse, error) {
	if s.FindByBrandOrModelFunc != nil {
		return s.FindByBrandOrModelFunc(ctx, input)
	}
	return nil, nil
}

func (s *MockSequencerService) Create(ctx context.Context, input models.SequencerCreateInput) (*models.SequencerAdminTableResponse, error) {
	if s.CreateFunc != nil {
		return s.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (s *MockSequencerService) Update(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error) {
	if s.UpdateFunc != nil {
		return s.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (s *MockSequencerService) Delete(ctx context.Context, ID uuid.UUID) error {
	if s.DeleteFunc != nil {
		return s.DeleteFunc(ctx, ID)
	}
	return nil
}
