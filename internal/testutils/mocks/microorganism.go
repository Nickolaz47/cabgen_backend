package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockMicroorganismRepository struct {
	GetMicroorganismsFunc          func(ctx context.Context) ([]models.Microorganism, error)
	GetActiveMicroorganismsFunc    func(ctx context.Context) ([]models.Microorganism, error)
	GetMicroorganismByIDFunc       func(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error)
	GetMicroorganismsBySpeciesFunc func(ctx context.Context, input, lang string) ([]models.Microorganism, error)
	GetMicroorganismDuplicateFunc  func(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error)
	CreateMicroorganismFunc        func(ctx context.Context, micro *models.Microorganism) error
	UpdateMicroorganismFunc        func(ctx context.Context, micro *models.Microorganism) error
	DeleteMicroorganismFunc        func(ctx context.Context, micro *models.Microorganism) error
}

func (r *MockMicroorganismRepository) GetMicroorganisms(ctx context.Context) ([]models.Microorganism, error) {
	if r.GetMicroorganismsFunc != nil {
		return r.GetMicroorganismsFunc(ctx)
	}
	return nil, nil
}

func (r *MockMicroorganismRepository) GetActiveMicroorganisms(ctx context.Context) ([]models.Microorganism, error) {
	if r.GetActiveMicroorganismsFunc != nil {
		return r.GetActiveMicroorganismsFunc(ctx)
	}
	return nil, nil
}

func (r *MockMicroorganismRepository) GetMicroorganismByID(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
	if r.GetMicroorganismByIDFunc != nil {
		return r.GetMicroorganismByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (r *MockMicroorganismRepository) GetMicroorganismsBySpecies(ctx context.Context, input, lang string) ([]models.Microorganism, error) {
	if r.GetMicroorganismsBySpeciesFunc != nil {
		return r.GetMicroorganismsBySpeciesFunc(ctx, input, lang)
	}
	return nil, nil
}

func (r *MockMicroorganismRepository) GetMicroorganismDuplicate(ctx context.Context, species string, variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
	if r.GetMicroorganismDuplicateFunc != nil {
		return r.GetMicroorganismDuplicateFunc(ctx, species, variety, ID)
	}
	return nil, nil
}

func (r *MockMicroorganismRepository) CreateMicroorganism(ctx context.Context, micro *models.Microorganism) error {
	if r.CreateMicroorganismFunc != nil {
		return r.CreateMicroorganismFunc(ctx, micro)
	}
	return nil
}

func (r *MockMicroorganismRepository) UpdateMicroorganism(ctx context.Context, micro *models.Microorganism) error {
	if r.UpdateMicroorganismFunc != nil {
		return r.UpdateMicroorganismFunc(ctx, micro)
	}
	return nil
}

func (r *MockMicroorganismRepository) DeleteMicroorganism(ctx context.Context, micro *models.Microorganism) error {
	if r.DeleteMicroorganismFunc != nil {
		return r.DeleteMicroorganismFunc(ctx, micro)
	}
	return nil
}

type MockMicroorganismService struct {
	FindAllFunc       func(ctx context.Context, language string) ([]models.MicroorganismAdminTableResponse, error)
	FindAllActiveFunc func(ctx context.Context, language string) ([]models.MicroorganismFormResponse, error)
	FindByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.MicroorganismAdminDetailResponse, error)
	FindBySpeciesFunc func(ctx context.Context, species, language string) ([]models.MicroorganismAdminTableResponse, error)
	CreateFunc        func(ctx context.Context, input models.MicroorganismCreateInput) (*models.MicroorganismAdminDetailResponse, error)
	UpdateFunc        func(ctx context.Context, ID uuid.UUID, input models.MicroorganismUpdateInput) (*models.MicroorganismAdminDetailResponse, error)
	DeleteFunc        func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockMicroorganismService) FindAll(ctx context.Context, language string) ([]models.MicroorganismAdminTableResponse, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, language)
	}
	return nil, nil
}

func (m *MockMicroorganismService) FindAllActive(ctx context.Context, language string) ([]models.MicroorganismFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx, language)
	}
	return nil, nil
}

func (m *MockMicroorganismService) FindByID(ctx context.Context, ID uuid.UUID) (*models.MicroorganismAdminDetailResponse, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockMicroorganismService) FindBySpecies(ctx context.Context, species, language string) ([]models.MicroorganismAdminTableResponse, error) {
	if m.FindBySpeciesFunc != nil {
		return m.FindBySpeciesFunc(ctx, species, language)
	}
	return nil, nil
}

func (m *MockMicroorganismService) Create(ctx context.Context, input models.MicroorganismCreateInput) (*models.MicroorganismAdminDetailResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockMicroorganismService) Update(ctx context.Context, ID uuid.UUID, input models.MicroorganismUpdateInput) (*models.MicroorganismAdminDetailResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (m *MockMicroorganismService) Delete(ctx context.Context, ID uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}
