package models

import "github.com/google/uuid"

type Sequencer struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Brand    string    `gorm:"not null" json:"brand"`
	Model    string    `gorm:"not null" json:"model"`
	IsActive bool      `gorm:"not null" json:"is_active"`
}

type SequencerFormResponse struct {
	ID    uuid.UUID `json:"id"`
	Model string    `json:"model"`
}

func (s *Sequencer) ToFormResponse() SequencerFormResponse {
	return SequencerFormResponse{
		ID:    s.ID,
		Model: s.Model,
	}
}

type SequencerCreateInput struct {
	Brand    string `json:"brand" binding:"required,min=3"`
	Model    string `json:"model" binding:"required,min=3"`
	IsActive bool   `json:"is_active"`
}

type SequencerUpdateInput struct {
	Brand    *string `json:"brand,omitempty" binding:"omitempty,min=3"`
	Model    *string `json:"model,omitempty" binding:"omitempty,min=3"`
	IsActive *bool   `json:"is_active,omitempty" binding:"omitempty"`
}
