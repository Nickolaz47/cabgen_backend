package models

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type HealthService struct {
	ID           string                   `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Name         string                   `gorm:"not null;uniqueIndex" json:"name"`
	Type         models.HealthServiceType `gorm:"not null" json:"type"`
	CountryID    uint                     `gorm:"not null" json:"-"`
	Country      models.Country           `gorm:"foreignKey:CountryID;references:ID"`
	City         string                   `gorm:"default:null" json:"city,omitempty"`
	Contactant   string                   `gorm:"default:null" json:"contactant,omitempty"`
	ContactEmail string                   `gorm:"default:null" json:"contact_email,omitempty"`
	ContactPhone string                   `gorm:"default:null" json:"contact_phone,omitempty"`
	IsActive     bool                     `gorm:"not null" json:"is_active"`
}

func NewHealthService(
	ID, name string, hServType models.HealthServiceType, country models.Country,
	city, contactant, contactEmail, contactPhone string, isActive bool,
) models.HealthService {
	return models.HealthService{
		ID:           uuid.MustParse(ID),
		Name:         name,
		Type:         hServType,
		CountryID:    country.ID,
		Country:      country,
		City:         city,
		Contactant:   contactant,
		ContactEmail: contactEmail,
		ContactPhone: contactPhone,
		IsActive:     isActive,
	}
}
