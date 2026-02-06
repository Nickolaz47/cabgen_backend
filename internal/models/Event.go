package models

import (
	"encoding/json"
	"time"
)

const (
	EventPending    = "PENDING"
	EventProcessing = "PROCESSING"
	EventDone       = "DONE"
	EventFailed     = "FAILED"
)

type Event struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null;index"`
	Payload   []byte `gorm:"not null"`
	Status    string `gorm:"not null;index"`
	Error     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func NewEvent(name string, data any) (*Event, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Event{Name: name, Payload: payload, Status: EventPending}, nil
}
