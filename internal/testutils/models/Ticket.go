package models

import (
	"time"

	rModels "github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Ticket struct {
	ID          uuid.UUID `gorm:"primaryKey;default:(hex(randomblob(16)))"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"not null"`
	Institution string    `gorm:"not null"`
	Subject     string    `gorm:"not null"`
	Message     string    `gorm:"not null"`
	Status      string    `gorm:"not null;default:'OPEN'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AdminID     *string       `gorm:"index"`
	Admin       *rModels.User `gorm:"foreignKey:AdminID;references:ID"`
}

func NewTicket(
	ID, name, email, institution, subject, message string,
	admin *rModels.User) rModels.Ticket {
	ticket := rModels.Ticket{
		ID:          uuid.MustParse(ID),
		Name:        name,
		Email:       email,
		Institution: institution,
		Subject:     subject,
		Message:     message,
		Status:      rModels.TicketStatusOpen,
	}

	if admin != nil {
		ticket.AdminID = &admin.ID
		ticket.Admin = admin
	}

	return ticket
}
