package models

import "github.com/google/uuid"

type Origin struct {
	ID       uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Names    map[string]string `gorm:"json;not null" json:"names"`
	IsActive bool              `gorm:"not null" json:"is_active"`
}
