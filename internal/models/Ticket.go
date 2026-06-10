package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	TicketStatusOpen       = "OPEN"
	TicketStatusInProgress = "IN_PROGRESS"
	TicketStatusResolved   = "RESOLVED"
)

type Ticket struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"not null"`
	Institution string    `gorm:"not null"`
	Subject     string    `gorm:"not null"`
	Message     string    `gorm:"not null"`
	Status      string    `gorm:"not null;default:'OPEN'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	AdminID *uuid.UUID `gorm:"type:uuid;index"`
	Admin   *User      `gorm:"foreignKey:AdminID;references:ID"`
}

type CreateTicketInput struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Institution string `json:"institution" binding:"required,max=150"`
	Subject     string `json:"subject" binding:"required,min=5,max=150"`
	Message     string `json:"message" binding:"required,min=10,max=2000"`
}

type TicketResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Institution string    `json:"institution"`
	Subject     string    `json:"subject"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`

	AdminID *uuid.UUID `json:"admin_id,omitempty"`
	Admin   string     `json:"admin,omitempty"`
}

func (t *Ticket) ToResponse() TicketResponse {
	response := TicketResponse{
		ID:          t.ID,
		Name:        t.Name,
		Email:       t.Email,
		Institution: t.Institution,
		Subject:     t.Subject,
		Message:     t.Message,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
	}

	if t.AdminID != nil {
		response.AdminID = t.AdminID
	}

	if t.Admin != nil {
		response.Admin = t.Admin.Username
	}

	return response
}
