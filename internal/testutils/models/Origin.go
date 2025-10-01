package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Origin struct {
	ID       string `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Pt       string `gorm:"not null" json:"pt"`
	En       string `gorm:"not null" json:"en"`
	Es       string `gorm:"not null" json:"es"`
	IsActive bool   `gorm:"not null" json:"is_active"`
}

func NewOrigin(ID, pt, en, es string, isActive bool) models.Origin {
	return models.Origin{
		ID:       uuid.MustParse(ID),
		Pt:       pt,
		En:       en,
		Es:       es,
		IsActive: isActive,
	}
}
