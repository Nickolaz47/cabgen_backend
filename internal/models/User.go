package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	Collaborator UserRole = "Collaborator"
	Admin        UserRole = "Admin"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Username    string     `gorm:"not null;uniqueIndex" json:"username"`
	Email       string     `gorm:"not null;uniqueIndex" json:"email"`
	Password    string     `gorm:"not null" json:"-"`
	CountryCode string     `gorm:"type:uuid;not null" json:"country_code"`
	Country     Country    `gorm:"foreignKey:CountryCode;references:Code"`
	IsActive    bool       `gorm:"not null" json:"is_active"`
	UserRole    UserRole   `gorm:"type:varchar(20);not null" json:"user_role"`
	Interest    *string    `gorm:"default:null" json:"interest,omitempty"`
	Role        *string    `gorm:"default:null" json:"role,omitempty"`
	Institution *string    `gorm:"default:null" json:"institution,omitempty"`
	CreatedBy   string     `gorm:"not null" json:"created_by"`
	ActivatedBy *string    `json:"activated_by"`
	ActivatedOn *time.Time `json:"activated_on"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Username:    u.Username,
		Email:       u.Email,
		CountryCode: u.CountryCode,
		Interest:    u.Interest,
		Role:        u.Role,
		Institution: u.Institution,
	}
}

type RegisterInput struct {
	Name            string `json:"name" binding:"required,min=3,max=100"`
	Username        string `json:"username" binding:"required,min=4,max=100"`
	Email           string `json:"email" binding:"required,email"`
	ConfirmEmail    string `json:"confirm_email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	CountryCode     string `json:"country_code" binding:"required"`

	Interest    *string `json:"interest,omitempty" binding:"omitempty,max=255"`
	Role        *string `json:"role,omitempty" binding:"omitempty,max=255"`
	Institution *string `json:"institution,omitempty" binding:"omitempty,max=255"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CountryCode string    `json:"country_code"`
	UserRole    string    `json:"user_role"`
	Interest    *string   `json:"interest,omitempty"`
	Role        *string   `json:"role,omitempty"`
	Institution *string   `json:"institution,omitempty"`
}
