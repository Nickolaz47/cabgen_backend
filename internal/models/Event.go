package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	EventPending    = "PENDING"
	EventProcessing = "PROCESSING"
	EventDone       = "DONE"
	EventFailed     = "FAILED"
)

type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null;index" json:"name"`
	Payload   []byte    `gorm:"not null" json:"payload"`
	Status    string    `gorm:"not null;index" json:"status"`
	Error     string    `gorm:"default:null" json:"error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewEvent(name string, data any) (*Event, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Event{Name: name, Payload: payload, Status: EventPending}, nil
}
