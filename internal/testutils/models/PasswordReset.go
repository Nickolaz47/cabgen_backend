package models

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type PasswordReset struct {
	ID        string    `gorm:"primaryKey;default:(hex(randomblob(16)))"`
	Email     string    `gorm:"not null;index"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func NewPasswordReset(email string, token string, expiresAt time.Time) models.PasswordReset {
	if email == "" {
		email = "user@example.com"
	}
	if token == "" {
		token = uuid.New().String()
	}
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(1 * time.Hour)
	}

	return models.PasswordReset{
		ID:        uuid.New(),
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}
