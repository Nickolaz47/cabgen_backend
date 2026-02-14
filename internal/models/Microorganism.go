package models

import (
	"github.com/google/uuid"
)

type Taxon string

const (
	Bacteria Taxon = "Bacteria"
	Virus    Taxon = "Virus"
	Protozoa Taxon = "Protozoa"
	Fungi    Taxon = "Fungi"
)

func (t Taxon) IsValid() bool {
	switch t {
	case Bacteria, Virus, Protozoa, Fungi:
		return true
	default:
		return false
	}
}

type Microorganism struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Taxon    Taxon     `gorm:"not null" json:"taxon"`
	Species  string    `gorm:"not null" json:"species"`
	Variety  JSONMap   `gorm:"type:jsonb" json:"variety"`
	IsActive bool      `gorm:"not null" json:"is_active"`
}

type MicroorganismAdminDetailResponse struct {
	ID       uuid.UUID         `json:"id"`
	Taxon    Taxon             `json:"taxon"`
	Species  string            `json:"species"`
	Variety  map[string]string `json:"variety"`
	IsActive bool              `json:"is_active"`
}

type MicroorganismAdminTableResponse struct {
	ID       uuid.UUID `json:"id"`
	Taxon    Taxon     `json:"taxon"`
	Species  string    `json:"species"`
	Variety  string    `json:"variety"`
	IsActive bool      `json:"is_active"`
}

type MicroorganismFormResponse struct {
	ID      uuid.UUID `json:"id"`
	Species string    `json:"species"`
}

func (m *Microorganism) ToAdminDetailResponse() MicroorganismAdminDetailResponse {
	return MicroorganismAdminDetailResponse{
		ID:       m.ID,
		Taxon:    m.Taxon,
		Species:  m.Species,
		Variety:  m.Variety,
		IsActive: m.IsActive,
	}
}

func (m *Microorganism) ToAdminTableResponse(
	language string) MicroorganismAdminTableResponse {
	if language == "" {
		language = "en"
	}

	return MicroorganismAdminTableResponse{
		ID:       m.ID,
		Taxon:    m.Taxon,
		Species:  m.Species,
		Variety:  m.Variety[language],
		IsActive: m.IsActive,
	}
}

func (m *Microorganism) ToFormResponse(language string) MicroorganismFormResponse {
	if language == "" {
		language = "en"
	}

	species := m.Species
	if m.Variety != nil {
		species = species + " " + m.Variety[language]
	}

	return MicroorganismFormResponse{
		ID:      m.ID,
		Species: species,
	}
}

type MicroorganismCreateInput struct {
	Taxon    Taxon             `json:"taxon" binding:"required"`
	Species  string            `json:"species" binding:"required,min=3"`
	Variety  map[string]string `json:"variety" binding:"min=3"`
	IsActive bool              `json:"is_active"`
}

type MicroorganismUpdateInput struct {
	Taxon    *Taxon            `json:"taxon,omitempty" binding:"omitempty"`
	Species  *string           `json:"species,omitempty" binding:"omitempty,min=3"`
	Variety  map[string]string `json:"variety,omitempty" binding:"omitempty,min=3"`
	IsActive *bool             `json:"is_active,omitempty" binding:"omitempty"`
}
