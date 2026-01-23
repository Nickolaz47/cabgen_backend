package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type SampleSource struct {
	ID       string            `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Names    map[string]string `gorm:"json;not null" json:"names"`
	Groups   map[string]string `gorm:"json;not null" json:"groups"`
	IsActive bool              `gorm:"not null" json:"is_active"`
}

func NewSampleSource(ID string, names, groups map[string]string, isActive bool) models.SampleSource {
	return models.SampleSource{
		ID:       uuid.MustParse(ID),
		Names:    names,
		Groups:   groups,
		IsActive: isActive,
	}
}
