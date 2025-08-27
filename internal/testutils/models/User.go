package models

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
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

var (
	name        = "Nicolas"
	password    = "12345678"
	email       = "nicolas@mail.com"
	countryCode = "BRA"
)

func NewLoginUser() models.User {
	return models.User{
		ID:          uuid.MustParse("e6d36ae2-4855-5bae-a76d-d29e3d57e76c"),
		Name:        "Nicolas",
		Username:    "nick",
		Password:    "$2a$10$P8SRTHBxlK09pYuj8Nn1A.2WMufAH1tZZKAPQel1bt0X5S82zbRGO",
		Email:       "nick@mail.com",
		CountryCode: "BRA",
		Country:     models.Country{Code: "BRA", Pt: "Brasil", Es: "Brazil", En: "Brazil"},
		IsActive:    true,
	}
}

func NewRegisterUser(username, inputEmail string) models.RegisterInput {
	if username == "" {
		username = "nmfaraujo"
	}

	if inputEmail == "" {
		inputEmail = email
	}

	return models.RegisterInput{
		Name:            name,
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
		Email:           inputEmail,
		ConfirmEmail:    inputEmail,
		CountryCode:     countryCode,
	}
}

func NewUpdateUserInput() models.UpdateUserInput {
	username := "nickol"
	interest := "Bacterial resistance"
	role := "Researcher"
	institution := "NCBI"

	return models.UpdateUserInput{
		Name:        &name,
		Username:    &username,
		CountryCode: &countryCode,
		Interest:    &interest,
		Role:        &role,
		Institution: &institution,
	}
}

func NewUserToken(id uuid.UUID, username string, userRole models.UserRole) models.UserToken {
	return models.UserToken{
		ID:       id,
		Username: username,
		UserRole: userRole,
	}
}
