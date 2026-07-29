package models

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type EmailUpdateRequest struct {
	ID        string    `gorm:"primaryKey;default:(hex(randomblob(16)))"`
	UserID    string    `gorm:"not null;index"`
	OldEmail  string    `gorm:"not null"`
	NewEmail  string    `gorm:"not null"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func NewEmailUpdateRequest(userID uuid.UUID, oldEmail, newEmail, token string,
	expiresAt time.Time) models.EmailUpdateRequest {
	if oldEmail == "" {
		oldEmail = "old@example.com"
	}
	if newEmail == "" {
		newEmail = "new@example.com"
	}
	if token == "" {
		token = uuid.New().String()
	}
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(1 * time.Hour)
	}

	return models.EmailUpdateRequest{
		ID:        uuid.New(),
		UserID:    userID,
		OldEmail:  oldEmail,
		NewEmail:  newEmail,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}
