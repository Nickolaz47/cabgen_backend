package models

import (
	"time"

	rModels "github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Audit struct {
	ID        uuid.UUID `gorm:"primaryKey;default:(hex(randomblob(16)))"`
	Event     string    `gorm:"not null"`
	Source    string    `gorm:"not null"`
	Status    int       `gorm:"not null"`
	Metadata  string
	CreatedAt time.Time
	UserID    *uuid.UUID    `gorm:"index"`
	User      *rModels.User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL"`
}

func NewAudit(event, source, metadata string, status int) rModels.Audit {
	if event == "" {
		event = rModels.AuditEventLogin
	}
	if source == "" {
		source = "127.0.0.1"
	}
	if metadata == "" {
		metadata = "{}"
	}
	if status == 0 {
		status = 200
	}

	admin := NewAdminLoginUser()
	return rModels.Audit{
		ID:        uuid.New(),
		Event:     event,
		Source:    source,
		Status:    status,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UserID:    &admin.ID,
		User:      &admin,
	}
}
