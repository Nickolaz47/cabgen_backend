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
	Password    string          `gorm:"not null" json:"-" `
	CountryID   uint            `gorm:"type:uint;not null" json:"country_id" `
	Country     models.Country  `gorm:"foreignKey:CountryID;references:ID" `
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
	name             = "Nicolas"
	username         = "nickol"
	password         = "12345678"
	email            = "nicolas@mail.com"
	countryID   uint = 1
	countryCode      = "BRA"
	interest         = "Bacterial resistance"
	role             = "Researcher"
	institution      = "NCBI"
)

func NewLoginUser() models.User {
	return models.User{
		ID:        uuid.MustParse("e6d36ae2-4855-5bae-a76d-d29e3d57e76c"),
		Name:      "Nicolas",
		Username:  "nick",
		Password:  "$2a$10$P8SRTHBxlK09pYuj8Nn1A.2WMufAH1tZZKAPQel1bt0X5S82zbRGO",
		Email:     "nick@mail.com",
		CountryID: countryID,
		Country: models.Country{
			ID:   countryID,
			Code: "BRA",
			Names: map[string]string{
				"pt": "Brasil",
				"en": "Brazil",
				"es": "Brazil",
			},
		},
		IsActive: true,
		UserRole: models.Collaborator,
	}
}

func NewAdminLoginUser() models.User {
	return models.User{
		ID:        uuid.MustParse("e6d36ae2-4855-5bae-a76d-d29e3d57e76d"),
		Name:      "Cabgen Admin",
		Username:  "admin",
		Password:  "$2a$10$P8SRTHBxlK09pYuj8Nn1A.2WMufAH1tZZKAPQel1bt0X5S82zbRGO",
		Email:     "admin@mail.com",
		UserRole:  models.Admin,
		CountryID: countryID,
		Country: models.Country{
			ID:   countryID,
			Code: "BRA",
			Names: map[string]string{
				"pt": "Brasil",
				"en": "Brazil",
				"es": "Brazil",
			},
		},
		IsActive:  true,
		CreatedBy: "system",
		CreatedAt: time.Now(),
	}
}

func NewRegisterUser(inputUsername, inputEmail string) models.UserRegisterInput {
	if inputUsername == "" {
		inputUsername = username
	}
	if inputEmail == "" {
		inputEmail = email
	}

	return models.UserRegisterInput{
		Name:            name,
		Username:        inputUsername,
		Email:           inputEmail,
		ConfirmEmail:    inputEmail,
		Password:        password,
		ConfirmPassword: password,
		CountryCode:     countryCode,
		Interest:        &interest,
		Role:            &role,
		Institution:     &institution,
	}
}

func NewUserUpdateInput() models.UserUpdateInput {
	return models.UserUpdateInput{
		Name:        &name,
		Username:    &username,
		CountryCode: &countryCode,
		Interest:    &interest,
		Role:        &role,
		Institution: &institution,
	}
}

func NewAdminCreateUserInput(inputEmail, inputUsername string) models.AdminUserCreateInput {
	if inputEmail == "" {
		inputEmail = "eddy@mail.com"
	}
	if inputUsername == "" {
		inputUsername = "eddy"
	}

	return models.AdminUserCreateInput{
		Name:        "Eddie",
		Username:    inputUsername,
		Email:       inputEmail,
		Password:    password,
		CountryCode: countryCode,
		UserRole:    models.Collaborator,
		IsActive:    true,
		Interest:    &interest,
		Role:        &role,
		Institution: &institution,
	}
}

func NewAdminUpdateUserInput() models.AdminUserUpdateInput {
	userRole := models.Collaborator
	isActive := true

	return models.AdminUserUpdateInput{
		Name:        &name,
		Username:    &username,
		Email:       &email,
		CountryCode: &countryCode,
		UserRole:    &userRole,
		IsActive:    &isActive,
		Interest:    &interest,
		Role:        &role,
		Institution: &institution,
	}
}

func NewUserToken(id uuid.UUID, username string, role models.UserRole) models.UserToken {
	return models.UserToken{
		ID:       id,
		Username: username,
		UserRole: role,
	}
}
