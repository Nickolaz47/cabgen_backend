package repositories

import (
	"context"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MicroorganismRepository interface {
	GetMicroorganisms(ctx context.Context) ([]models.Microorganism, error)
	GetActiveMicroorganisms(ctx context.Context) ([]models.Microorganism, error)
	GetMicroorganismByID(ctx context.Context,
		ID uuid.UUID) (*models.Microorganism, error)
	GetMicroorganismsBySpecies(ctx context.Context,
		input, lang string) ([]models.Microorganism, error)
	GetMicroorganismDuplicate(ctx context.Context, species string,
		variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error)
	CreateMicroorganism(ctx context.Context, micro *models.Microorganism) error
	UpdateMicroorganism(ctx context.Context, micro *models.Microorganism) error
	DeleteMicroorganism(ctx context.Context, micro *models.Microorganism) error
}

type microorganismRepo struct {
	DB *gorm.DB
}

func NewMicroorganismRepository(db *gorm.DB) MicroorganismRepository {
	return &microorganismRepo{
		DB: db,
	}
}

func (r *microorganismRepo) GetMicroorganisms(
	ctx context.Context,
) ([]models.Microorganism, error) {
	var microorganisms []models.Microorganism
	if err := r.DB.WithContext(ctx).Find(&microorganisms).Error; err != nil {
		return nil, err
	}

	return microorganisms, nil
}

func (r *microorganismRepo) GetActiveMicroorganisms(ctx context.Context) ([]models.Microorganism, error) {
	var microorganisms []models.Microorganism
	if err := r.DB.WithContext(ctx).Where("is_active = true").
		Find(&microorganisms).Error; err != nil {
		return nil, err
	}

	return microorganisms, nil
}

func (r *microorganismRepo) GetMicroorganismByID(ctx context.Context, ID uuid.UUID) (*models.Microorganism, error) {
	var microorganism models.Microorganism
	if err := r.DB.WithContext(ctx).Where("id = ?", ID).
		First(&microorganism).Error; err != nil {
		return nil, err
	}

	return &microorganism, nil
}

func (r *microorganismRepo) GetMicroorganismsBySpecies(ctx context.Context, input, lang string) ([]models.Microorganism, error) {
	var microorganisms []models.Microorganism
	query := "LOWER(species) LIKE LOWER(?) OR LOWER(variety->>'" + lang + "') LIKE LOWER(?)"
	if err := r.DB.WithContext(ctx).Where(query, "%"+input+"%", "%"+input+"%").
		Find(&microorganisms).Error; err != nil {
		return nil, err
	}

	return microorganisms, nil
}

func (r *microorganismRepo) GetMicroorganismDuplicate(
	ctx context.Context, species string,
	variety models.JSONMap, ID uuid.UUID) (*models.Microorganism, error) {
	var microorganism models.Microorganism

	conditions := r.DB.WithContext(ctx)
	if len(variety) != 0 {
		for lang, value := range variety {
			conditions = conditions.Or(
				fmt.Sprintf(
					"LOWER(species) = LOWER(?) AND LOWER(variety->>'%s') = LOWER(?)",
					lang,
				), species, value,
			)
		}
	} else {
		conditions = conditions.Where("LOWER(species) = LOWER(?)", species)
	}

	query := r.DB.WithContext(ctx).Where(conditions)

	if ID != uuid.Nil {
		query = query.Where("id != ?", ID)
	}

	if err := query.First(&microorganism).Error; err != nil {
		return nil, err
	}

	return &microorganism, nil
}

func (r *microorganismRepo) CreateMicroorganism(ctx context.Context, micro *models.Microorganism) error {
	return r.DB.WithContext(ctx).Create(micro).Error
}

func (r *microorganismRepo) UpdateMicroorganism(ctx context.Context, micro *models.Microorganism) error {
	return r.DB.WithContext(ctx).Save(micro).Error
}

func (r *microorganismRepo) DeleteMicroorganism(ctx context.Context, micro *models.Microorganism) error {
	return r.DB.WithContext(ctx).Delete(micro).Error
}
