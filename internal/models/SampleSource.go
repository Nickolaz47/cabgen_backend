package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SampleSource struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Names    JSONMap   `gorm:"type:jsonb;not null" json:"names"`
	Groups   JSONMap   `gorm:"type:jsonb;not null" json:"groups"`
	IsActive bool      `gorm:"not null" json:"is_active"`
}

type SampleSourceResponse struct {
	Name     string `json:"name"`
	Group    string `json:"group"`
	IsActive bool   `json:"is_active"`
}

type SampleSourceFormResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (s *SampleSource) ToResponse(c *gin.Context) SampleSourceResponse {
	language := c.GetHeader("Accept-Language")
	if language == "" {
		language = "en"
	}

	name := s.Names[language]
	group := s.Groups[language]

	return SampleSourceResponse{
		Name:     name,
		Group:    group,
		IsActive: s.IsActive,
	}
}

func (s *SampleSource) ToFormResponse(c *gin.Context) SampleSourceFormResponse {
	language := c.GetHeader("Accept-Language")
	if language == "" {
		language = "en"
	}

	name := s.Names[language]

	return SampleSourceFormResponse{
		ID:   s.ID,
		Name: name,
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
