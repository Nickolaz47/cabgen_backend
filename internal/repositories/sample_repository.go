package repositories

import (
	"context"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SampleRepository interface {
	GetSamples(ctx context.Context, input string) ([]models.Sample, error)
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

func (s *sampleRepo) GetSamples(ctx context.Context,
	input string) ([]models.Sample, error) {
	var samples []models.Sample

	query := s.DB.WithContext(ctx)
	// Name or Run number or Origin code or Microorganism species
	if input != "" {
		searchTerm := "%" + strings.ToLower(input) + "%"

		query = query.Preload("Microorganism").
			Joins("JOIN microorganisms ON microorganisms.id"+
				" = samples.microorganism_id").Where(
			`
			LOWER(samples.name) LIKE ? OR 
			LOWER(samples.run_number) LIKE ? OR 
			LOWER(samples.origin_code) LIKE ? OR
			LOWER(microorganisms.species) LIKE ?
			`,
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	if err := query.Find(&samples).Error; err != nil {
		return nil, err
	}

	return samples, nil
}

func (s *sampleRepo) GetSampleByID(ctx context.Context,
	ID uuid.UUID) (*models.Sample, error) {
	var sample models.Sample
	if err := s.DB.WithContext(ctx).Where("id = ?", ID).
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
