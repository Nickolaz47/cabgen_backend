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
	return m.GetHealthServicesFunc(ctx)
}

func (m *MockHealthServiceRepository) GetActiveHealthServices(ctx context.Context) ([]models.HealthService, error) {
	return m.GetActiveHealthServicesFunc(ctx)
}

func (m *MockHealthServiceRepository) GetHealthServiceByID(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
	return m.GetHealthServiceByIDFunc(ctx, ID)
}

func (m *MockHealthServiceRepository) GetHealthServicesByName(ctx context.Context, input string) ([]models.HealthService, error) {
	return m.GetHealthServicesByNameFunc(ctx, input)
}

func (m *MockHealthServiceRepository) GetHealthServiceDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
	return m.GetHealthServiceDuplicateFunc(ctx, name, ID)
}

func (m *MockHealthServiceRepository) CreateHealthService(ctx context.Context, healthService *models.HealthService) error {
	return m.CreateHealthServiceFunc(ctx, healthService)
}

func (m *MockHealthServiceRepository) UpdateHealthService(ctx context.Context, healthService *models.HealthService) error {
	return m.UpdateHealthServiceFunc(ctx, healthService)
}

func (m *MockHealthServiceRepository) DeleteHealthService(ctx context.Context, healthService *models.HealthService) error {
	return m.DeleteHealthServiceFunc(ctx, healthService)
}
