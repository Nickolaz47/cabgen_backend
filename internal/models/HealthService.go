package models

import "github.com/google/uuid"

type HealthServiceType string

const (
	Public  HealthServiceType = "Public"
	Private HealthServiceType = "Private"
)

func (h HealthServiceType) IsValid() bool {
	switch h {
	case Public, Private:
		return true
	default:
		return false
	}
}

var HealthServiceTypes = []HealthServiceType{Public, Private}

type HealthService struct {
	ID           uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string            `gorm:"type:varchar(255);not null;uniqueIndex"`
	Type         HealthServiceType `gorm:"type:varchar(20);not null"`
	CountryID    uint              `gorm:"not null"`
	Country      Country           `gorm:"foreignKey:CountryID;references:ID"`
	City         *string           `gorm:"type:varchar(255);default:null"`
	Contactant   *string           `gorm:"type:varchar(255);default:null"`
	ContactEmail *string           `gorm:"type:varchar(255);default:null"`
	ContactPhone *string           `gorm:"type:varchar(255);default:null"`
	IsActive     bool              `gorm:"not null"`
}

type HealthServiceAdminTableResponse struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	Type         HealthServiceType `json:"type"`
	Country      string            `json:"country"`
	City         *string           `json:"city,omitempty"`
	Contactant   *string           `json:"contactant,omitempty"`
	ContactEmail *string           `json:"contact_email,omitempty"`
	ContactPhone *string           `json:"contact_phone,omitempty"`
	IsActive     bool              `json:"is_active"`
}

type HealthServiceFormResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (h *HealthService) ToAdminTableResponse() HealthServiceAdminTableResponse {
	return HealthServiceAdminTableResponse{
		ID:           h.ID,
		Name:         h.Name,
		Type:         h.Type,
		Country:      h.Country.Code,
		City:         h.City,
		Contactant:   h.Contactant,
		ContactEmail: h.ContactEmail,
		ContactPhone: h.ContactPhone,
		IsActive:     h.IsActive,
	}
}

func (h *HealthService) ToFormResponse() HealthServiceFormResponse {
	return HealthServiceFormResponse{
		ID:   h.ID,
		Name: h.Name,
	}
}

type HealthServiceCreateInput struct {
	Name         string            `json:"name" binding:"required,min=3"`
	Type         HealthServiceType `json:"type" binding:"required,min=3"`
	CountryCode  string            `json:"country_code" binding:"required,len=3"`
	City         *string           `json:"city,omitempty" binding:"omitempty,min=3"`
	Contactant   *string           `json:"contactant,omitempty" binding:"omitempty,min=3"`
	ContactEmail *string           `json:"contact_email,omitempty" binding:"omitempty,email"`
	ContactPhone *string           `json:"contact_phone,omitempty" binding:"omitempty,e164"`
	IsActive     bool              `json:"is_active"`
}

type HealthServiceUpdateInput struct {
	Name         *string            `json:"name,omitempty" binding:"omitempty,min=3"`
	Type         *HealthServiceType `json:"type,omitempty" binding:"omitempty,min=3"`
	CountryCode  *string            `json:"country_code,omitempty" binding:"omitempty,len=3"`
	City         *string            `json:"city,omitempty" binding:"omitempty,min=3"`
	Contactant   *string            `json:"contactant,omitempty" binding:"omitempty,min=3"`
	ContactEmail *string            `json:"contact_email,omitempty" binding:"omitempty,email"`
	ContactPhone *string            `json:"contact_phone,omitempty" binding:"omitempty,e164"`
	IsActive     *bool              `json:"is_active,omitempty" binding:"omitempty"`
}
