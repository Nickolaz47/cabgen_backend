package repository

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SampleSourceRepository struct {
	DB *gorm.DB
}

func NewSampleSourceRepo(db *gorm.DB) *SampleSourceRepository {
	return &SampleSourceRepository{DB: db}
}

func (r *SampleSourceRepository) GetSampleSources() ([]models.SampleSource, error) {
	var sampleSources []models.SampleSource
	if err := r.DB.Find(&sampleSources).Error; err != nil {
		return nil, err
	}

	return sampleSources, nil
}

func (r *SampleSourceRepository) GetActiveSampleSources() ([]models.SampleSource, error) {
	var sampleSources []models.SampleSource
	if err := r.DB.Where("is_active = true").Find(&sampleSources).Error; err != nil {
		return nil, err
	}

	return sampleSources, nil
}

func (r *SampleSourceRepository) GetSampleSourceByID(ID uuid.UUID) (*models.SampleSource, error) {
	var sampleSource models.SampleSource
	if err := r.DB.Where("id = ?", ID).First(&sampleSource).Error; err != nil {
		return nil, err
	}

	return &sampleSource, nil
}

func (r *SampleSourceRepository) GetSampleSourcesByNameOrGroup(input, lang string) ([]models.SampleSource, error) {
	var sampleSources []models.SampleSource
	query := "LOWER(names->>'" + lang + "') LIKE LOWER(?) OR LOWER(groups->>'" + lang + "') LIKE LOWER(?)"
	if err := r.DB.Where(query, "%"+input+"%", "%"+input+"%").Find(&sampleSources).Error; err != nil {
		return nil, err
	}

	return sampleSources, nil
}

func (r *SampleSourceRepository) CreateSampleSource(sampleSource *models.SampleSource) error {
	return r.DB.Create(sampleSource).Error
}

func (r *SampleSourceRepository) UpdateSampleSource(sampleSource *models.SampleSource) error {
	return r.DB.Save(sampleSource).Error
}

func (r *SampleSourceRepository) DeleteSampleSource(sampleSource *models.SampleSource) error {
	return r.DB.Delete(sampleSource).Error
}
