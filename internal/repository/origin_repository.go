package repository

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OriginRepository struct {
	DB *gorm.DB
}

func NewOriginRepo(db *gorm.DB) *OriginRepository {
	return &OriginRepository{DB: db}
}

func (r *OriginRepository) GetOrigins() ([]models.Origin, error) {
	var origins []models.Origin
	if err := r.DB.Find(&origins).Error; err != nil {
		return nil, err
	}

	return origins, nil
}

func (r *OriginRepository) GetOriginByID(ID uuid.UUID) (*models.Origin, error) {
	var origin models.Origin
	if err := r.DB.Where("id = ?", ID).First(&origin).Error; err != nil {
		return nil, err
	}

	return &origin, nil
}

func (r *OriginRepository) GetOriginByName(name string) (*models.Origin, error) {
	var origin models.Origin
	if err := r.DB.Where("pt = ? OR en = ? OR es = ?", name, name, name).First(&origin).Error; err != nil {
		return nil, err
	}

	return &origin, nil
}

func (r *OriginRepository) CreateOrigin(origin *models.Origin) error {
	return r.DB.Create(origin).Error
}

func (r *OriginRepository) UpdateOrigin(origin *models.Origin) error {
	return r.DB.Save(origin).Error
}

func (r *OriginRepository) DeleteOrigin(origin *models.Origin) error {
	return r.DB.Delete(origin).Error
}
