package repository

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
