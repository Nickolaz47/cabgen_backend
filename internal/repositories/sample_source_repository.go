package repositories

import (
	"context"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SampleSourceRepository interface {
	GetSampleSources(ctx context.Context) ([]models.SampleSource, error)
	GetActiveSampleSources(ctx context.Context) ([]models.SampleSource, error)
	GetSampleSourceByID(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error)
	GetSampleSourcesByNameOrGroup(ctx context.Context, input, lang string) ([]models.SampleSource, error)
	GetSampleSourceDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error)
	CreateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error
	UpdateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error
	DeleteSampleSource(ctx context.Context, sampleSource *models.SampleSource) error
}

type sampleSourceRepo struct {
	DB *gorm.DB
}

func NewSampleSourceRepo(db *gorm.DB) SampleSourceRepository {
	return &sampleSourceRepo{DB: db}
}

func (r *sampleSourceRepo) GetSampleSources(ctx context.Context) ([]models.SampleSource, error) {
	var sampleSources []models.SampleSource
	if err := r.DB.WithContext(ctx).Find(&sampleSources).Error; err != nil {
		return nil, err
	}

	return sampleSources, nil
}

func (r *sampleSourceRepo) GetActiveSampleSources(ctx context.Context) ([]models.SampleSource, error) {
	var sampleSources []models.SampleSource
	if err := r.DB.WithContext(ctx).Where("is_active = true").Find(&sampleSources).Error; err != nil {
		return nil, err
	}

	return sampleSources, nil
}

func (r *sampleSourceRepo) GetSampleSourceByID(ctx context.Context, ID uuid.UUID) (*models.SampleSource, error) {
	var sampleSource models.SampleSource
	if err := r.DB.WithContext(ctx).Where("id = ?", ID).First(&sampleSource).Error; err != nil {
		return nil, err
	}

	return &sampleSource, nil
}

func (r *sampleSourceRepo) GetSampleSourcesByNameOrGroup(ctx context.Context, input, lang string) ([]models.SampleSource, error) {
	var sampleSources []models.SampleSource
	query := "LOWER(names->>'" + lang + "') LIKE LOWER(?) OR LOWER(groups->>'" + lang + "') LIKE LOWER(?)"
	if err := r.DB.WithContext(ctx).Where(query, "%"+input+"%", "%"+input+"%").Find(&sampleSources).Error; err != nil {
		return nil, err
	}

	return sampleSources, nil
}

func (r *sampleSourceRepo) GetSampleSourceDuplicate(ctx context.Context, names models.JSONMap, ID uuid.UUID) (*models.SampleSource, error) {
	var sampleSource models.SampleSource

	query := r.DB.WithContext(ctx)

	for lang, value := range names {
		query = query.Or(
			fmt.Sprintf(
				"LOWER(names->>'%s') LIKE LOWER(?)",
				lang,
			),
			value,
		)
	}

	if ID != uuid.Nil {
		query = query.Where("id != ?", ID)
	}

	if err := query.First(&sampleSource).Error; err != nil {
		return nil, err
	}

	return &sampleSource, nil
}

func (r *sampleSourceRepo) CreateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	return r.DB.WithContext(ctx).Create(sampleSource).Error
}

func (r *sampleSourceRepo) UpdateSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	return r.DB.WithContext(ctx).Save(sampleSource).Error
}

func (r *sampleSourceRepo) DeleteSampleSource(ctx context.Context, sampleSource *models.SampleSource) error {
	return r.DB.WithContext(ctx).Delete(sampleSource).Error
}
