package repository

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SequencerRepository interface {
	GetSequencers(ctx context.Context) ([]models.Sequencer, error)
	GetActiveSequencers(ctx context.Context) ([]models.Sequencer, error)
	GetSequencerByID(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error)
	GetSequencersByBrandOrModel(ctx context.Context, input string) ([]models.Sequencer, error)
	GetSequencerDuplicate(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error)
	CreateSequencer(ctx context.Context, sequencer *models.Sequencer) error
	UpdateSequencer(ctx context.Context, sequencer *models.Sequencer) error
	DeleteSequencer(ctx context.Context, sequencer *models.Sequencer) error
}

type sequencerRepo struct {
	DB *gorm.DB
}

func NewSequencerRepo(db *gorm.DB) SequencerRepository {
	return &sequencerRepo{DB: db}
}

func (r *sequencerRepo) GetSequencers(ctx context.Context) ([]models.Sequencer, error) {
	var sequencers []models.Sequencer
	if err := r.DB.WithContext(ctx).Find(&sequencers).Error; err != nil {
		return nil, err
	}

	return sequencers, nil
}

func (r *sequencerRepo) GetActiveSequencers(ctx context.Context) ([]models.Sequencer, error) {
	var sequencers []models.Sequencer
	if err := r.DB.WithContext(ctx).Where("is_active = true").Find(&sequencers).Error; err != nil {
		return nil, err
	}

	return sequencers, nil
}

func (r *sequencerRepo) GetSequencerByID(ctx context.Context, ID uuid.UUID) (*models.Sequencer, error) {
	var sequencer models.Sequencer
	if err := r.DB.WithContext(ctx).Where("id = ?", ID).First(&sequencer).Error; err != nil {
		return nil, err
	}

	return &sequencer, nil
}

func (r *sequencerRepo) GetSequencersByBrandOrModel(ctx context.Context, input string) ([]models.Sequencer, error) {
	var sequencers []models.Sequencer
	inputQuery := "%" + input + "%"
	if err := r.DB.WithContext(ctx).Where("LOWER(model) LIKE LOWER(?) OR LOWER(brand) LIKE LOWER(?)", inputQuery, inputQuery).Find(&sequencers).Error; err != nil {
		return nil, err
	}

	return sequencers, nil
}

func (r *sequencerRepo) GetSequencerDuplicate(ctx context.Context, model string, ID uuid.UUID) (*models.Sequencer, error) {
	var sequencer models.Sequencer

	query := r.DB.WithContext(ctx).Where("LOWER(model) = LOWER(?)", model)

	if ID != uuid.Nil {
		query = query.Where("id != ?", ID)
	}

	if err := query.First(&sequencer).Error; err != nil {
		return nil, err
	}

	return &sequencer, nil
}

func (r *sequencerRepo) CreateSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	return r.DB.WithContext(ctx).Create(sequencer).Error
}

func (r *sequencerRepo) UpdateSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	return r.DB.WithContext(ctx).Save(sequencer).Error
}

func (r *sequencerRepo) DeleteSequencer(ctx context.Context, sequencer *models.Sequencer) error {
	return r.DB.WithContext(ctx).Delete(sequencer).Error
}
