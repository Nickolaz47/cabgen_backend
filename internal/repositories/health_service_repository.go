package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HealthServiceRepository interface {
	GetHealthServices(ctx context.Context) ([]models.HealthService, error)
	GetActiveHealthServices(ctx context.Context) ([]models.HealthService, error)
	GetHealthServiceByID(ctx context.Context, ID uuid.UUID) (*models.HealthService, error)
	GetHealthServicesByName(ctx context.Context, input string) ([]models.HealthService, error)
	GetHealthServiceDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error)
	CreateHealthService(ctx context.Context, healthService *models.HealthService) error
	UpdateHealthService(ctx context.Context, healthService *models.HealthService) error
	DeleteHealthService(ctx context.Context, healthService *models.HealthService) error
}

type healthServiceRepo struct {
	DB *gorm.DB
}

func NewHealthServiceRepo(db *gorm.DB) HealthServiceRepository {
	return &healthServiceRepo{DB: db}
}

func (r *healthServiceRepo) GetHealthServices(ctx context.Context) ([]models.HealthService, error) {
	var healthServices []models.HealthService

	if err := r.DB.WithContext(ctx).Preload("Country").Find(&healthServices).Error; err != nil {
		return nil, err
	}

	return healthServices, nil
}

func (r *healthServiceRepo) GetActiveHealthServices(ctx context.Context) ([]models.HealthService, error) {
	var healthServices []models.HealthService

	if err := r.DB.WithContext(ctx).Preload("Country").Where("is_active = true").Find(&healthServices).Error; err != nil {
		return nil, err
	}

	return healthServices, nil
}

func (r *healthServiceRepo) GetHealthServiceByID(ctx context.Context, ID uuid.UUID) (*models.HealthService, error) {
	var healthService models.HealthService

	if err := r.DB.WithContext(ctx).Preload("Country").Where("id = ?", ID).First(&healthService).Error; err != nil {
		return nil, err
	}

	return &healthService, nil
}

func (r *healthServiceRepo) GetHealthServicesByName(ctx context.Context, input string) ([]models.HealthService, error) {
	var healthServices []models.HealthService
	inputQuery := "%" + input + "%"

	if err := r.DB.WithContext(ctx).Preload("Country").Where("LOWER(name) LIKE LOWER(?)", inputQuery).Find(&healthServices).Error; err != nil {
		return nil, err
	}

	return healthServices, nil
}

func (r *healthServiceRepo) GetHealthServiceDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.HealthService, error) {
	var healthService models.HealthService

	query := r.DB.WithContext(ctx).Preload("Country").Where("LOWER(name) = LOWER(?)", name)

	if ID != uuid.Nil {
		query = query.Where("id != ?", ID)
	}

	if err := query.First(&healthService).Error; err != nil {
		return nil, err
	}

	return &healthService, nil
}

func (r *healthServiceRepo) CreateHealthService(ctx context.Context, healthService *models.HealthService) error {
	return r.DB.WithContext(ctx).Create(healthService).Error
}

func (r *healthServiceRepo) UpdateHealthService(ctx context.Context, healthService *models.HealthService) error {
	return r.DB.WithContext(ctx).Save(healthService).Error
}

func (r *healthServiceRepo) DeleteHealthService(ctx context.Context, healthService *models.HealthService) error {
	return r.DB.WithContext(ctx).Delete(healthService).Error
}
