package models

import (
	"github.com/google/uuid"
)

type SampleSource struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Names    JSONMap   `gorm:"type:jsonb;not null" json:"names"`
	Groups   JSONMap   `gorm:"type:jsonb;not null" json:"groups"`
	IsActive bool      `gorm:"not null" json:"is_active"`
}

type SampleSourceAdminDetailResponse struct {
	ID       uuid.UUID         `json:"id"`
	Names    map[string]string `json:"names"`
	Groups   map[string]string `json:"groups"`
	IsActive bool              `json:"is_active"`
}

type SampleSourceAdminTableResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Group    string    `json:"group"`
	IsActive bool      `json:"is_active"`
}

type SampleSourceFormResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Group string    `json:"group"`
}

func (s *SampleSource) ToAdminDetailResponse() SampleSourceAdminDetailResponse {
	return SampleSourceAdminDetailResponse{
		ID:       s.ID,
		Names:    s.Names,
		Groups:   s.Groups,
		IsActive: s.IsActive,
	}
}

func (s *SampleSource) ToAdminTableResponse(language string) SampleSourceAdminTableResponse {
	if language == "" {
		language = "en"
	}

	return SampleSourceAdminTableResponse{
		ID:       s.ID,
		Name:     s.Names[language],
		Group:    s.Groups[language],
		IsActive: s.IsActive,
	}
}

func (s *SampleSource) ToFormResponse(language string) SampleSourceFormResponse {
	if language == "" {
		language = "en"
	}

	name := s.Names[language]
	group := s.Groups[language]

	return SampleSourceFormResponse{
		ID:    s.ID,
		Name:  name,
		Group: group,
	}
}

type SampleSourceCreateInput struct {
	Names    map[string]string `json:"names" binding:"required,min=3"`
	Groups   map[string]string `json:"groups" binding:"required,min=3"`
	IsActive bool              `json:"is_active"`
}

type SampleSourceUpdateInput struct {
	Names    map[string]string `json:"names,omitempty" binding:"omitempty,min=3"`
	Groups   map[string]string `json:"groups,omitempty" binding:"omitempty,min=3"`
	IsActive *bool             `json:"is_active,omitempty" binding:"omitempty"`
}
