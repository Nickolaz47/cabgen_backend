package models

import (
	"github.com/google/uuid"
)

type Origin struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Names    JSONMap   `gorm:"type:jsonb;not null" json:"names"`
	IsActive bool      `gorm:"not null" json:"is_active"`
}

type OriginAdminDetailResponse struct {
	ID       uuid.UUID         `json:"id"`
	Names    map[string]string `json:"names"`
	IsActive bool              `json:"is_active"`
}

type OriginAdminTableResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
}

type OriginFormResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (o *Origin) ToAdminDetailResponse() OriginAdminDetailResponse {
	return OriginAdminDetailResponse{
		ID:       o.ID,
		Names:    o.Names,
		IsActive: o.IsActive,
	}
}

func (o *Origin) ToAdminTableResponse(language string) OriginAdminTableResponse {
	if language == "" {
		language = "en"
	}

	return OriginAdminTableResponse{
		ID:       o.ID,
		Name:     o.Names[language],
		IsActive: o.IsActive,
	}
}

func (o *Origin) ToFormResponse(language string) OriginFormResponse {
	if language == "" {
		language = "en"
	}

	return OriginFormResponse{
		ID:   o.ID,
		Name: o.Names[language],
	}
}

type OriginCreateInput struct {
	Names    map[string]string `json:"names" binding:"required,min=3"`
	IsActive bool              `json:"is_active"`
}

type OriginUpdateInput struct {
	Names    map[string]string `json:"names,omitempty" binding:"omitempty,min=3"`
	IsActive *bool             `json:"is_active,omitempty" binding:"omitempty"`
}
