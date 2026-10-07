package repositories

import (
	"context"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SampleRepository interface {
	GetSamples(ctx context.Context, input string,
		userID uuid.UUID, limit, offset int) ([]models.Sample, int64,
		error)
	GetSampleByID(ctx context.Context, ID uuid.UUID) (*models.Sample, error)
	CreateSample(ctx context.Context, sample *models.Sample) error
	UpdateSample(ctx context.Context, sample *models.Sample) error
	DeleteSample(ctx context.Context, sample *models.Sample) error
}

type sampleRepo struct {
	DB *gorm.DB
}

func NewSampleRepo(db *gorm.DB) SampleRepository {
	return &sampleRepo{DB: db}
}

func (s *sampleRepo) GetSamples(ctx context.Context, input string,
	userID uuid.UUID, limit, offset int) ([]models.Sample, int64, error) {
	var samples []models.Sample
	var total int64

	build := func() *gorm.DB {
		query := s.DB.WithContext(ctx).Model(&models.Sample{})
		// Filter by Origin code
		if input != "" {
			searchTerm := "%" + strings.ToLower(input) + "%"
			query = query.Where("LOWER(samples.origin_code) LIKE ?",
				searchTerm)
		}

		if userID != uuid.Nil {
			query = query.Where("samples.user_id = ?", userID)
		}

		return query
	}

	if limit > 0 {
		if err := build().Count(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	query := build().
		Preload("Country").
		Preload("User").
		Preload("Origin").
		Preload("SampleSource").
		Preload("Microorganism").
		Preload("Sequencer").
		Preload("Laboratory").
		Preload("HealthService")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("samples.created_at DESC, samples.id")

	if err := query.Find(&samples).Error; err != nil {
		return nil, 0, err
	}

	return samples, total, nil
}

func (s *sampleRepo) GetSampleByID(ctx context.Context,
	ID uuid.UUID) (*models.Sample, error) {
	var sample models.Sample
	if err := s.DB.WithContext(ctx).
		Preload("Country").
		Preload("User").
		Preload("Origin").
		Preload("SampleSource").
		Preload("Microorganism").
		Preload("Sequencer").
		Preload("Laboratory").
		Preload("HealthService").
		Where("id = ?", ID).
		First(&sample).Error; err != nil {
		return nil, err
	}

	return &sample, nil
}

func (s *sampleRepo) CreateSample(ctx context.Context,
	sample *models.Sample) error {
	return s.DB.WithContext(ctx).Create(sample).Error
}

func (s *sampleRepo) UpdateSample(ctx context.Context,
	sample *models.Sample) error {
	return s.DB.WithContext(ctx).Save(sample).Error
}

func (s *sampleRepo) DeleteSample(ctx context.Context,
	sample *models.Sample) error {
	return s.DB.WithContext(ctx).Delete(sample).Error
}
