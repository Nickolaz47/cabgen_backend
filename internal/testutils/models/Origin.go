package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Origin struct {
	ID       string            `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Names    map[string]string `gorm:"json;not null" json:"names"`
	IsActive bool              `gorm:"not null" json:"is_active"`
}

func NewOrigin(ID string, names map[string]string, isActive bool) models.Origin {
	return models.Origin{
		ID:       uuid.MustParse(ID),
		Names:    names,
		IsActive: isActive,
	}
}
