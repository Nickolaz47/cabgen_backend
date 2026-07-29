package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailUpdateRequest struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	OldEmail  string    `gorm:"not null"`
	NewEmail  string    `gorm:"not null"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (r *EmailUpdateRequest) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

type RequestEmailUpdateInput struct {
	NewEmail        string `json:"new_email" binding:"required,email"`
	ConfirmNewEmail string `json:"confirm_new_email" binding:"required,eqfield=NewEmail"`
}

type ConfirmEmailUpdateInput struct {
	Token string `json:"token" binding:"required"`
}
