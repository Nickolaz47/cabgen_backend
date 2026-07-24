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

type OriginSeedRepository struct {
	DB *gorm.DB
}

func NewOriginSeedRepository(db *gorm.DB) *OriginSeedRepository {
	return &OriginSeedRepository{DB: db}
}

func (r *OriginSeedRepository) BulkInsert(ctx context.Context,
	origins []models.Origin) error {
	return r.DB.WithContext(ctx).Create(&origins).Error
}

func (r *OriginSeedRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Origin{}).
		Count(&count).Error
	return count, err
}

type SequencerSeedRepository struct {
	DB *gorm.DB
}

func NewSequencerSeedRepository(db *gorm.DB) *SequencerSeedRepository {
	return &SequencerSeedRepository{DB: db}
}

func (r *SequencerSeedRepository) BulkInsert(ctx context.Context,
	sequencers []models.Sequencer) error {
	return r.DB.WithContext(ctx).Create(&sequencers).Error
}

func (r *SequencerSeedRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Sequencer{}).
		Count(&count).Error
	return count, err
}

type LaboratorySeedRepository struct {
	DB *gorm.DB
}

func NewLaboratorySeedRepository(db *gorm.DB) *LaboratorySeedRepository {
	return &LaboratorySeedRepository{DB: db}
}

func (r *LaboratorySeedRepository) BulkInsert(ctx context.Context,
	laboratories []models.Laboratory) error {
	return r.DB.WithContext(ctx).Create(&laboratories).Error
}

func (r *LaboratorySeedRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Laboratory{}).
		Count(&count).Error
	return count, err
}

type SampleSourceSeedRepository struct {
	DB *gorm.DB
}

func NewSampleSourceSeedRepository(db *gorm.DB) *SampleSourceSeedRepository {
	return &SampleSourceSeedRepository{DB: db}
}

func (r *SampleSourceSeedRepository) BulkInsert(ctx context.Context,
	sampleSources []models.SampleSource) error {
	return r.DB.WithContext(ctx).Create(&sampleSources).Error
}

func (r *SampleSourceSeedRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.SampleSource{}).
		Count(&count).Error
	return count, err
}
