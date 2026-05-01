package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockHealthServiceRepository struct {
	GetHealthServicesFunc         func(ctx context.Context) ([]models.HealthService, error)
	GetActiveHealthServicesFunc   func(ctx context.Context) ([]models.HealthService, error)
	GetHealthServiceByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.HealthService, error)
	GetHealthServicesByNameFunc   func(ctx context.Context, input string) ([]models.HealthService, error)
	GetHealthServiceDuplicateFunc func(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error)
	CreateHealthServiceFunc       func(ctx context.Context, healthService *models.HealthService) error
	UpdateHealthServiceFunc       func(ctx context.Context, healthService *models.HealthService) error
	DeleteHealthServiceFunc       func(ctx context.Context, healthService *models.HealthService) error
}

func (m *MockHealthServiceRepository) GetHealthServices(ctx context.Context) ([]models.HealthService, error) {
	if m.GetHealthServicesFunc != nil {
		return m.GetHealthServicesFunc(ctx)
	}

	return nil, nil
}

func (m *MockHealthServiceRepository) GetActiveHealthServices(ctx context.Context) ([]models.HealthService, error) {
	if m.GetActiveHealthServicesFunc != nil {
		return m.GetActiveHealthServicesFunc(ctx)
	}

	return nil, nil
}

func (m *MockHealthServiceRepository) GetHealthServiceByID(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
	if m.GetHealthServiceByIDFunc != nil {
		return m.GetHealthServiceByIDFunc(ctx, ID)
	}

	return nil, nil
}

func (m *MockHealthServiceRepository) GetHealthServicesByName(ctx context.Context, input string) ([]models.HealthService, error) {
	if m.GetHealthServicesByNameFunc != nil {
		return m.GetHealthServicesByNameFunc(ctx, input)
	}

	return nil, nil
}

func (m *MockHealthServiceRepository) GetHealthServiceDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
	if m.GetHealthServiceDuplicateFunc != nil {
		return m.GetHealthServiceDuplicateFunc(ctx, name, ID)
	}

	return nil, nil
}

func (m *MockHealthServiceRepository) CreateHealthService(ctx context.Context, healthService *models.HealthService) error {
	if m.CreateHealthServiceFunc != nil {
		return m.CreateHealthServiceFunc(ctx, healthService)
	}

	return nil
}

func (m *MockHealthServiceRepository) UpdateHealthService(ctx context.Context, healthService *models.HealthService) error {
	if m.UpdateHealthServiceFunc != nil {
		return m.UpdateHealthServiceFunc(ctx, healthService)
	}

	return nil
}

func (m *MockHealthServiceRepository) DeleteHealthService(ctx context.Context, healthService *models.HealthService) error {
	if m.DeleteHealthServiceFunc != nil {
		return m.DeleteHealthServiceFunc(ctx, healthService)
	}

	return nil
}

type MockHealthServiceService struct {
	FindAllFunc       func(ctx context.Context) ([]models.HealthServiceAdminTableResponse, error)
	FindAllActiveFunc func(ctx context.Context) ([]models.HealthServiceFormResponse, error)
	FindByIDFunc      func(ctx context.Context, ID uuid.UUID) (*models.HealthServiceAdminTableResponse, error)
	FindByNameFunc    func(ctx context.Context, name string) ([]models.HealthServiceAdminTableResponse, error)
	CreateFunc        func(ctx context.Context, input models.HealthServiceCreateInput) (*models.HealthServiceAdminTableResponse, error)
	UpdateFunc        func(ctx context.Context, ID uuid.UUID, input models.HealthServiceUpdateInput) (*models.HealthServiceAdminTableResponse, error)
	DeleteFunc        func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockHealthServiceService) FindAll(ctx context.Context) ([]models.HealthServiceAdminTableResponse, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}

	return nil, nil
}

func (m *MockHealthServiceService) FindAllActive(ctx context.Context) ([]models.HealthServiceFormResponse, error) {
	if m.FindAllActiveFunc != nil {
		return m.FindAllActiveFunc(ctx)
	}

	return nil, nil
}

func (m *MockHealthServiceService) FindByID(ctx context.Context, ID uuid.UUID) (*models.HealthServiceAdminTableResponse, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}

	return nil, nil
}

func (m *MockHealthServiceService) FindByName(ctx context.Context, name string) ([]models.HealthServiceAdminTableResponse, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(ctx, name)
	}

	return nil, nil
}

func (m *MockHealthServiceService) Create(ctx context.Context, input models.HealthServiceCreateInput) (*models.HealthServiceAdminTableResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}

	return nil, nil
}

func (m *MockHealthServiceService) Update(ctx context.Context, ID uuid.UUID, input models.HealthServiceUpdateInput) (*models.HealthServiceAdminTableResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input)
	}

	return nil, nil
}

func (m *MockHealthServiceService) Delete(ctx context.Context, ID uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}

	return nil
}
