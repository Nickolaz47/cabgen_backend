package models

import "github.com/google/uuid"

type Origin struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Pt       string    `gorm:"not null" json:"pt"`
	En       string    `gorm:"not null" json:"en"`
	Es       string    `gorm:"not null" json:"es"`
	IsActive bool      `gorm:"not null" json:"is_active"`
}
