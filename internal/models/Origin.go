package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Origin struct {
	ID       uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Names    map[string]string `gorm:"json;not null" json:"names"`
	IsActive bool              `gorm:"not null" json:"is_active"`
}

type OriginResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
}

func (o *Origin) ToResponse(c *gin.Context) OriginResponse {
	language := c.GetHeader("Accept-Language")
	if language == "" {
		language = "en"
	}

	name := o.Names[language]

	return OriginResponse{
		ID:       o.ID,
		Name:     name,
		IsActive: o.IsActive,
	}
}

type OriginCreateInput struct {
	Names    map[string]string `json:"names" binding:"required,min=3"`
	IsActive bool              `json:"is_active" binding:"required"`
}

type OriginUpdateInput struct {
	Names    map[string]string `json:"names,omitempty" binding:"omitempty"`
	IsActive *bool             `json:"is_active,omitempty" binding:"omitempty"`
}
