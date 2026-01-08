package models

import "github.com/google/uuid"

type Sequencer struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Model    string    `gorm:"not null" json:"model"`
	Brand    string    `gorm:"not null" json:"brand"`
	IsActive bool      `gorm:"not null" json:"is_active"`
}

type SequencerAdminTableResponse struct {
	ID       uuid.UUID `json:"id"`
	Model    string    `json:"model"`
	Brand    string    `json:"brand"`
	IsActive bool      `json:"is_active"`
}

type SequencerFormResponse struct {
	ID    uuid.UUID `json:"id"`
	Model string    `json:"model"`
	Brand string    `json:"brand"`
}

func (s *Sequencer) ToAdminTableResponse() SequencerAdminTableResponse {
	return SequencerAdminTableResponse{
		ID:       s.ID,
		Model:    s.Model,
		Brand:    s.Brand,
		IsActive: s.IsActive,
	}
}

func (s *Sequencer) ToFormResponse() SequencerFormResponse {
	return SequencerFormResponse{
		ID:    s.ID,
		Model: s.Model,
		Brand: s.Brand,
	}
}

type SequencerCreateInput struct {
	Model    string `json:"model" binding:"required,min=3"`
	Brand    string `json:"brand" binding:"required,min=3"`
	IsActive bool   `json:"is_active"`
}

type SequencerUpdateInput struct {
	Model    *string `json:"model,omitempty" binding:"omitempty,min=3"`
	Brand    *string `json:"brand,omitempty" binding:"omitempty,min=3"`
	IsActive *bool   `json:"is_active,omitempty" binding:"omitempty"`
}
