package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Sequencer struct {
	ID       string `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Model    string `gorm:"not null" json:"model"`
	Brand    string `gorm:"not null" json:"brand"`
	IsActive bool   `gorm:"not null" json:"is_active"`
}

func NewSequencer(ID, model, brand string, isActive bool) models.Sequencer {
	return models.Sequencer{
		ID:       uuid.MustParse(ID),
		Model:    model,
		Brand:    brand,
		IsActive: isActive,
	}
}
