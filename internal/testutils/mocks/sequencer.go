package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockSequencerRepository struct {
	GetSequencersFunc               func(ctx context.Context) ([]models.Sequencer, error)
	GetActiveSequencersFunc         func(ctx context.Context) ([]models.Sequencer, error)
	GetSequencerByIDFunc            func(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error)
	GetSequencersByBrandOrModelFunc func(ctx context.Context, input string) ([]models.Sequencer, error)
	GetSequencerDuplicateFunc       func(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error)
	CreateSequencerFunc             func(ctx context.Context, sequencer *models.Sequencer) error
	UpdateSequencerFunc             func(ctx context.Context, sequencer *models.Sequencer) error
	DeleteSequencerFunc             func(ctx context.Context, sequencer *models.Sequencer) error
}

func (s *MockSequencerRepository) GetSequencers(ctx context.Context) ([]models.Sequencer, error) {
	if s.GetSequencersFunc != nil {
		return s.GetSequencersFunc(ctx)
	}
	return nil, nil
}

func (s *MockSequencerRepository) GetActiveSequencers(ctx context.Context) ([]models.Sequencer, error) {
	if s.GetActiveSequencersFunc != nil {
		return s.GetActiveSequencersFunc(ctx)
	}
	return nil, nil
}

func (s *MockSequencerRepository) GetSequencerByID(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
	if s.GetSequencerByIDFunc != nil {
		return s.GetSequencerByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (s *MockSequencerRepository) GetSequencersByBrandOrModel(ctx context.Context, input string) ([]models.Sequencer, error) {
	if s.GetSequencersByBrandOrModelFunc != nil {
		return s.GetSequencersByBrandOrModelFunc(ctx, input)
	}
	return nil, nil
}

func (s *MockSequencerRepository) GetSequencerDuplicate(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
	if s.GetSequencerDuplicateFunc != nil {
		return s.GetSequencerDuplicateFunc(ctx, model, ID)
	}
	return nil, nil
}

func (s *MockSequencerRepository) CreateSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	if s.CreateSequencerFunc != nil {
		return s.CreateSequencerFunc(ctx, sequencer)
	}
	return nil
}

func (s *MockSequencerRepository) UpdateSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	if s.UpdateSequencerFunc != nil {
		return s.UpdateSequencerFunc(ctx, sequencer)
	}
	return nil
}

func (s *MockSequencerRepository) DeleteSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	if s.DeleteSequencerFunc != nil {
		return s.DeleteSequencerFunc(ctx, sequencer)
	}
	return nil
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
