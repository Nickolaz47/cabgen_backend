package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/gorm"
)

type CountrySeedRepository struct {
	DB *gorm.DB
}

func NewCountrySeedRepository(db *gorm.DB) *CountrySeedRepository {
	return &CountrySeedRepository{DB: db}
}

func (r *CountrySeedRepository) BulkInsert(ctx context.Context, countries []models.Country) error {
	return r.DB.WithContext(ctx).Create(&countries).Error
}

func (r *CountrySeedRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Country{}).Count(&count).Error
	return count, err
}

type MicroorganismSeedRepository struct {
	DB *gorm.DB
}

func NewMicroorganismSeedRepository(db *gorm.DB) *MicroorganismSeedRepository {
	return &MicroorganismSeedRepository{DB: db}
}

func (r *MicroorganismSeedRepository) BulkInsert(ctx context.Context,
	microorganisms []models.Microorganism) error {
	return r.DB.WithContext(ctx).Create(&microorganisms).Error
}

func (r *MicroorganismSeedRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Microorganism{}).
		Count(&count).Error
	return count, err
}
