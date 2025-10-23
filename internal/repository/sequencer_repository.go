package repository

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SequencerRepository struct {
	DB *gorm.DB
}

func NewSequencerRepo(db *gorm.DB) *SequencerRepository {
	return &SequencerRepository{DB: db}
}

func (r *SequencerRepository) GetSequencers() ([]models.Sequencer, error) {
	var sequencers []models.Sequencer
	if err := r.DB.Find(&sequencers).Error; err != nil {
		return nil, err
	}

	return sequencers, nil
}

func (r *SequencerRepository) GetActiveSequencers() ([]models.Sequencer, error) {
	var sequencers []models.Sequencer
	if err := r.DB.Where("is_active = true").Find(&sequencers).Error; err != nil {
		return nil, err
	}

	return sequencers, nil
}

func (r *SequencerRepository) GetSequencerByID(ID uuid.UUID) (*models.Sequencer, error) {
	var sequencer models.Sequencer
	if err := r.DB.Where("id = ?", ID).First(&sequencer).Error; err != nil {
		return nil, err
	}

	return &sequencer, nil
}

func (r *SequencerRepository) GetSequencersByBrandOrModel(input string) ([]models.Sequencer, error) {
	var sequencers []models.Sequencer
	inputQuery := "%" + input + "%"
	if err := r.DB.Where("LOWER(model) LIKE LOWER(?) OR LOWER(brand) LIKE LOWER(?)", inputQuery, inputQuery).Find(&sequencers).Error; err != nil {
		return nil, err
	}

	return sequencers, nil
}

func (r *SequencerRepository) CreateSequencer(sequencer *models.Sequencer) error {
	return r.DB.Create(sequencer).Error
}

func (r *SequencerRepository) UpdateSequencer(sequencer *models.Sequencer) error {
	return r.DB.Save(sequencer).Error
}

func (r *SequencerRepository) DeleteSequencer(sequencer *models.Sequencer) error {
	return r.DB.Delete(sequencer).Error
}
