package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Laboratory struct {
	ID           string `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Name         string `gorm:"not null" json:"name"`
	Abbreviation string `gorm:"not null" json:"abbreviation"`
	IsActive     bool   `gorm:"not null" json:"is_active"`
}

func NewLaboratory(ID, name, abbreviation string, isActive bool) models.Laboratory {
	return models.Laboratory{
		ID:           uuid.MustParse(ID),
		Name:         name,
		Abbreviation: abbreviation,
		IsActive:     isActive,
	}
}
