package models

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

type User struct {
	ID          string          `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Username    string          `gorm:"not null;uniqueIndex" json:"username"`
	Email       string          `gorm:"not null;uniqueIndex" json:"email"`
	Password    string          `gorm:"not null" json:"-"`
	CountryCode string          `gorm:"type:string;not null" json:"country_code"`
	Country     models.Country  `gorm:"foreignKey:CountryCode;references:Code"`
	IsActive    bool            `gorm:"not null" json:"is_active"`
	UserRole    models.UserRole `gorm:"type:varchar(20);not null" json:"user_role"`
	Interest    *string         `gorm:"default:null" json:"interest,omitempty"`
	Role        *string         `gorm:"default:null" json:"role,omitempty"`
	Institution *string         `gorm:"default:null" json:"institution,omitempty"`
	CreatedBy   string          `gorm:"not null" json:"created_by"`
	ActivatedBy *string         `json:"activated_by"`
	ActivatedOn *time.Time      `json:"activated_on"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
