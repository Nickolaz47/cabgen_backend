package repositories

import (
	"context"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OriginRepository interface {
	GetOrigins(ctx context.Context) ([]models.Origin, error)
	GetActiveOrigins(ctx context.Context) ([]models.Origin, error)
	GetOriginByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error)
	GetOriginsByName(ctx context.Context, name, lang string) ([]models.Origin, error)
	GetOriginDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error)
	CreateOrigin(ctx context.Context, origin *models.Origin) error
	UpdateOrigin(ctx context.Context, origin *models.Origin) error
	DeleteOrigin(ctx context.Context, origin *models.Origin) error
}

type originRepo struct {
	DB *gorm.DB
}

func NewOriginRepo(db *gorm.DB) OriginRepository {
	return &originRepo{DB: db}
}

func (r *originRepo) GetOrigins(ctx context.Context) ([]models.Origin, error) {
	var origins []models.Origin
	if err := r.DB.WithContext(ctx).Find(&origins).Error; err != nil {
		return nil, err
	}

	return origins, nil
}

func (r *originRepo) GetActiveOrigins(ctx context.Context) ([]models.Origin, error) {
	var origins []models.Origin
	if err := r.DB.WithContext(ctx).Where("is_active = true").Find(&origins).Error; err != nil {
		return nil, err
	}

	return origins, nil
}

func (r *originRepo) GetOriginByID(ctx context.Context, ID uuid.UUID) (*models.Origin, error) {
	var origin models.Origin
	if err := r.DB.WithContext(ctx).Where("id = ?", ID).First(&origin).Error; err != nil {
		return nil, err
	}

	return &origin, nil
}

func (r *originRepo) GetOriginsByName(ctx context.Context, name, lang string) ([]models.Origin, error) {
	var origins []models.Origin
	query := "LOWER(names->>'" + lang + "') LIKE LOWER(?)"
	if err := r.DB.WithContext(ctx).Where(query, "%"+name+"%").Find(&origins).Error; err != nil {
		return nil, err
	}

	return origins, nil
}

func (r *originRepo) GetOriginDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.Origin, error) {
	var origin models.Origin

	conditions := r.DB.WithContext(ctx)
	for lang, value := range names {
		conditions = conditions.Or(
			fmt.Sprintf(
				"LOWER(names->>'%s') = LOWER(?)",
				lang,
			),
			value,
		)
	}

	query := r.DB.WithContext(ctx).Where(conditions)

	if ID != uuid.Nil {
		query = query.Where("id != ?", ID)
	}

	if err := query.First(&origin).Error; err != nil {
		return nil, err
	}

	return &origin, nil
}

func (r *originRepo) CreateOrigin(ctx context.Context, origin *models.Origin) error {
	return r.DB.WithContext(ctx).Create(origin).Error
}

func (r *originRepo) UpdateOrigin(ctx context.Context, origin *models.Origin) error {
	return r.DB.WithContext(ctx).Save(origin).Error
}

func (r *originRepo) DeleteOrigin(ctx context.Context, origin *models.Origin) error {
	return r.DB.WithContext(ctx).Delete(origin).Error
}
