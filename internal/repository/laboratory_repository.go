package repository

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LaboratoryRepository interface {
	GetLaboratories(ctx context.Context) ([]models.Laboratory, error)
	GetActiveLaboratories(ctx context.Context) ([]models.Laboratory, error)
	GetLaboratoryByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error)
	GetLaboratoriesByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error)
	GetLaboratoryDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error)
	CreateLaboratory(ctx context.Context, lab *models.Laboratory) error
	UpdateLaboratory(ctx context.Context, lab *models.Laboratory) error
	DeleteLaboratory(ctx context.Context, lab *models.Laboratory) error
}

type laboratoryRepo struct {
	DB *gorm.DB
}

func NewLaboratoryRepo(db *gorm.DB) LaboratoryRepository {
	return &laboratoryRepo{DB: db}
}

func (r *laboratoryRepo) GetLaboratories(ctx context.Context) ([]models.Laboratory, error) {
	var labs []models.Laboratory

	if err := r.DB.WithContext(ctx).Find(&labs).Error; err != nil {
		return nil, err
	}

	return labs, nil
}

func (r *laboratoryRepo) GetActiveLaboratories(ctx context.Context) ([]models.Laboratory, error) {
	var labs []models.Laboratory
	if err := r.DB.WithContext(ctx).Where("is_active = true").Find(&labs).Error; err != nil {
		return nil, err
	}

	return labs, nil
}

func (r *laboratoryRepo) GetLaboratoryByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
	var lab models.Laboratory
	if err := r.DB.WithContext(ctx).Where("id = ?", ID).First(&lab).Error; err != nil {
		return nil, err
	}

	return &lab, nil
}

func (r *laboratoryRepo) GetLaboratoriesByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error) {
	var labs []models.Laboratory
	inputQuery := "%" + input + "%"
	if err := r.DB.WithContext(ctx).Where("LOWER(name) LIKE LOWER(?) OR LOWER(abbreviation) LIKE LOWER(?)", inputQuery, inputQuery).Find(&labs).Error; err != nil {
		return nil, err
	}

	return labs, nil
}

func (r *laboratoryRepo) GetLaboratoryDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error) {
	var lab models.Laboratory

	query := r.DB.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name)

	if ID != uuid.Nil {
		query = query.Where("id != ?", ID)
	}

	if err := query.First(&lab).Error; err != nil {
		return nil, err
	}

	return &lab, nil
}

func (r *laboratoryRepo) CreateLaboratory(ctx context.Context, lab *models.Laboratory) error {
	return r.DB.WithContext(ctx).Create(lab).Error
}

func (r *laboratoryRepo) UpdateLaboratory(ctx context.Context, lab *models.Laboratory) error {
	return r.DB.WithContext(ctx).Save(lab).Error
}

func (r *laboratoryRepo) DeleteLaboratory(ctx context.Context, lab *models.Laboratory) error {
	return r.DB.WithContext(ctx).Delete(lab).Error
}
