package models

import "github.com/google/uuid"

type Laboratory struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Abbreviation string    `gorm:"not null" json:"abbreviation"`
	IsActive     bool      `gorm:"not null" json:"is_active"`
}

type LaboratoryFormResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (s *Laboratory) ToFormResponse() LaboratoryFormResponse {
	return LaboratoryFormResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}

type LaboratoryCreateInput struct {
	Name         string `json:"name" binding:"required,min=3"`
	Abbreviation string `json:"abbreviation" binding:"required,min=3"`
	IsActive     bool   `json:"is_active"`
}

type LaboratoryUpdateInput struct {
	Name         *string `json:"name,omitempty" binding:"omitempty,min=3"`
	Abbreviation *string `json:"abbreviation,omitempty" binding:"omitempty,min=3"`
	IsActive     *bool   `json:"is_active,omitempty" binding:"omitempty"`
}
