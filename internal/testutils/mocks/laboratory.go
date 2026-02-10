package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockLaboratoryRepository struct {
	GetLaboratoriesFunc                     func(ctx context.Context) ([]models.Laboratory, error)
	GetActiveLaboratoriesFunc               func(ctx context.Context) ([]models.Laboratory, error)
	GetLaboratoryByIDFunc                   func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error)
	GetLaboratoriesByNameOrAbbreviationFunc func(ctx context.Context, input string) ([]models.Laboratory, error)
	GetLaboratoryDuplicateFunc              func(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error)
	CreateLaboratoryFunc                    func(ctx context.Context, lab *models.Laboratory) error
	UpdateLaboratoryFunc                    func(ctx context.Context, lab *models.Laboratory) error
	DeleteLaboratoryFunc                    func(ctx context.Context, lab *models.Laboratory) error
}

func (r *MockLaboratoryRepository) GetLaboratories(ctx context.Context) ([]models.Laboratory, error) {
	if r.GetLaboratoriesFunc != nil {
		return r.GetLaboratoriesFunc(ctx)
	}

	return nil, nil
}

func (r *MockLaboratoryRepository) GetActiveLaboratories(ctx context.Context) ([]models.Laboratory, error) {
	if r.GetActiveLaboratoriesFunc != nil {
		return r.GetActiveLaboratoriesFunc(ctx)
	}

	return nil, nil
}

func (r *MockLaboratoryRepository) GetLaboratoryByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
	if r.GetLaboratoryByIDFunc != nil {
		return r.GetLaboratoryByIDFunc(ctx, ID)
	}

	return nil, nil
}

func (r *MockLaboratoryRepository) GetLaboratoriesByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error) {
	if r.GetLaboratoriesByNameOrAbbreviationFunc != nil {
		return r.GetLaboratoriesByNameOrAbbreviationFunc(ctx, input)
	}

	return nil, nil
}

func (r *MockLaboratoryRepository) GetLaboratoryDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
	if r.GetLaboratoryDuplicateFunc != nil {
		return r.GetLaboratoryDuplicateFunc(ctx, name, ID)
	}

	return nil, nil
}

func (r *MockLaboratoryRepository) CreateLaboratory(ctx context.Context, lab *models.Laboratory) error {
	if r.CreateLaboratoryFunc != nil {
		return r.CreateLaboratoryFunc(ctx, lab)
	}

	return nil
}

func (r *MockLaboratoryRepository) UpdateLaboratory(ctx context.Context, lab *models.Laboratory) error {
	if r.UpdateLaboratoryFunc != nil {
		return r.UpdateLaboratoryFunc(ctx, lab)
	}

	return nil
}

func (r *MockLaboratoryRepository) DeleteLaboratory(ctx context.Context, lab *models.Laboratory) error {
	if r.DeleteLaboratoryFunc != nil {
		return r.DeleteLaboratoryFunc(ctx, lab)
	}

	return nil
}

type MockLaboratoryService struct {
	FindAllFunc                  func(ctx context.Context) ([]models.LaboratoryAdminTableResponse, error)
	FindAllActiveFunc            func(ctx context.Context) ([]models.LaboratoryFormResponse, error)
	FindByIDFunc                 func(ctx context.Context, ID uuid.UUID) (*models.LaboratoryAdminTableResponse, error)
	FindByNameOrAbbreviationFunc func(ctx context.Context, input string) ([]models.LaboratoryAdminTableResponse, error)
	CreateFunc                   func(ctx context.Context, input models.LaboratoryCreateInput) (*models.LaboratoryAdminTableResponse, error)
	UpdateFunc                   func(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.LaboratoryAdminTableResponse, error)
	DeleteFunc                   func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockLaboratoryService) FindAll(
	ctx context.Context,
) ([]models.LaboratoryAdminTableResponse, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindAllActive(
	ctx context.Context,
) ([]models.LaboratoryFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindByID(
	ctx context.Context,
	ID uuid.UUID,
) (*models.LaboratoryAdminTableResponse, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (m *MockLaboratoryService) FindByNameOrAbbreviation(
	ctx context.Context,
	input string,
) ([]models.LaboratoryAdminTableResponse, error) {
	if m.FindByNameOrAbbreviationFunc != nil {
		return m.FindByNameOrAbbreviationFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockLaboratoryService) Create(
	ctx context.Context,
	input models.LaboratoryCreateInput,
) (*models.LaboratoryAdminTableResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockLaboratoryService) Update(
	ctx context.Context,
	ID uuid.UUID,
	input models.LaboratoryUpdateInput,
) (*models.LaboratoryAdminTableResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}
	return nil, nil
}

func (m *MockLaboratoryService) Delete(
	ctx context.Context,
	ID uuid.UUID,
) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}
