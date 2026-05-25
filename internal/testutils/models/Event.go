package models

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Event struct {
	ID        string    `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Name      string    `gorm:"not null;index" json:"name"`
	Payload   []byte    `gorm:"not null" json:"payload"`
	Status    string    `gorm:"not null;index" json:"status"`
	Error     string    `gorm:"default:null" json:"error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewEvent(ID, name string, payload []byte,
	status, err string) models.Event {
	return models.Event{
		ID:      uuid.MustParse(ID),
		Name:    name,
		Payload: payload,
		Status:  status,
		Error:   err,
	}
}
