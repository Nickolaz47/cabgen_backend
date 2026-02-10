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

func (r UserRole) IsValid() bool {
	switch r {
	case Collaborator, Admin:
		return true
	default:
		return false
	}
}

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Username    string     `gorm:"not null;uniqueIndex" json:"username"`
	Email       string     `gorm:"not null;uniqueIndex" json:"email"`
	Password    string     `gorm:"not null" json:"-"`
	CountryID   uint       `gorm:"not null" json:"-"`
	Country     Country    `gorm:"foreignKey:CountryID;references:ID"`
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

type UserResponse struct {
	Name        string   `json:"name"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	CountryCode string   `json:"country_code"`
	Country     string   `json:"country"`
	UserRole    UserRole `json:"user_role"`
	Interest    *string  `json:"interest,omitempty"`
	Role        *string  `json:"role,omitempty"`
	Institution *string  `json:"institution,omitempty"`
}

type AdminUserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	CountryCode string     `json:"country_code"`
	Country     string     `json:"country"`
	UserRole    UserRole   `json:"user_role"`
	IsActive    bool       `json:"is_active"`
	CreatedBy   string     `json:"created_by"`
	ActivatedBy *string    `json:"activated_by"`
	ActivatedOn *time.Time `json:"activated_on"`
	Interest    *string    `json:"interest,omitempty"`
	Role        *string    `json:"role,omitempty"`
	Institution *string    `json:"institution,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *User) ToResponse(language string) UserResponse {
	if language == "" {
		language = "en"
	}

	countryName := u.Country.Names[language]
	return UserResponse{
		Name:        u.Name,
		Username:    u.Username,
		Email:       u.Email,
		CountryCode: u.Country.Code,
		Country:     countryName,
		UserRole:    u.UserRole,
		Interest:    u.Interest,
		Role:        u.Role,
		Institution: u.Institution,
	}
}

func (u *User) ToAdminResponse(language string) AdminUserResponse {
	if language == "" {
		language = "en"
	}

	countryName := u.Country.Names[language]
	return AdminUserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Username:    u.Username,
		Email:       u.Email,
		CountryCode: u.Country.Code,
		Country:     countryName,
		UserRole:    u.UserRole,
		IsActive:    u.IsActive,
		CreatedBy:   u.CreatedBy,
		ActivatedBy: u.ActivatedBy,
		ActivatedOn: u.ActivatedOn,
		Interest:    u.Interest,
		Role:        u.Role,
		Institution: u.Institution,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func (u *User) ToToken() UserToken {
	return UserToken{
		ID:       u.ID,
		Username: u.Username,
		UserRole: u.UserRole,
	}
}

type UserRegisterInput struct {
	Name            string `json:"name" binding:"required,min=3,max=100"`
	Username        string `json:"username" binding:"required,min=4,max=100"`
	Email           string `json:"email" binding:"required,email"`
	ConfirmEmail    string `json:"confirm_email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	CountryCode     string `json:"country_code" binding:"required,len=3"`

	Interest    *string `json:"interest,omitempty" binding:"omitempty,max=255"`
	Role        *string `json:"role,omitempty" binding:"omitempty,max=255"`
	Institution *string `json:"institution,omitempty" binding:"omitempty,max=255"`
}

type UserUpdateInput struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Username    *string `json:"username,omitempty" binding:"omitempty,min=4,max=100"`
	CountryCode *string `json:"country_code,omitempty" binding:"omitempty,len=3"`

	Interest    *string `json:"interest,omitempty" binding:"omitempty,max=255"`
	Role        *string `json:"role,omitempty" binding:"omitempty,max=255"`
	Institution *string `json:"institution,omitempty" binding:"omitempty,max=255"`
}

type AdminUserCreateInput struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Username    string `json:"username" binding:"required,min=4,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=32"`
	CountryCode string `json:"country_code" binding:"required,len=3"`

	UserRole UserRole `json:"user_role" binding:"required"`
	IsActive bool     `json:"is_active"`

	Interest    *string `json:"interest,omitempty" binding:"omitempty,max=255"`
	Role        *string `json:"role,omitempty" binding:"omitempty,max=255"`
	Institution *string `json:"institution,omitempty" binding:"omitempty,max=255"`
}

type AdminUserUpdateInput struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Username    *string `json:"username,omitempty" binding:"omitempty,min=4,max=100"`
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	Password    *string `json:"password" binding:"omitempty,min=8,max=32"`
	CountryCode *string `json:"country_code,omitempty" binding:"omitempty,len=3"`

	UserRole *UserRole `json:"user_role,omitempty"`
	IsActive *bool     `json:"is_active,omitempty"`

	Interest    *string `json:"interest,omitempty" binding:"omitempty,max=255"`
	Role        *string `json:"role,omitempty" binding:"omitempty,max=255"`
	Institution *string `json:"institution,omitempty" binding:"omitempty,max=255"`
}

type AdminUserFilter struct {
	// Name or username or email
	Input    *string
	UserRole *UserRole
	Active   *bool
}
