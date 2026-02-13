package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Microorganism struct {
	ID       string            `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Taxon    models.Taxon      `gorm:"not null" json:"taxon"`
	Species  string            `gorm:"not null" json:"species"`
	Variety  map[string]string `gorm:"json" json:"variety"`
	IsActive bool              `gorm:"not null" json:"is_active"`
}

func NewMicroorganism(
	ID string, taxon models.Taxon,
	species string, variety map[string]string,
	isActive bool) models.Microorganism {
	return models.Microorganism{
		ID:       uuid.MustParse(ID),
		Taxon:    taxon,
		Species:  species,
		Variety:  variety,
		IsActive: isActive,
	}
}
