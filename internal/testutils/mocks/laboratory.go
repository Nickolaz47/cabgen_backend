package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

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
