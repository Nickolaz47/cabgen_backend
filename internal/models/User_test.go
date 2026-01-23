package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestUserToResponse(t *testing.T) {
	mockUser := testmodels.NewLoginUser()
	lang := "en"

	mockResponse := mockUser.ToResponse(lang)
	expected := models.UserResponse{
		Name:        mockUser.Name,
		Username:    mockUser.Username,
		Email:       mockUser.Email,
		CountryCode: mockUser.Country.Code,
		Country:     mockUser.Country.Names[lang],
		UserRole:    mockUser.UserRole,
		Interest:    mockUser.Interest,
		Role:        mockUser.Role,
		Institution: mockUser.Institution,
	}

	assert.Equal(t, expected, mockResponse)
}

func TestUserToAdminResponse(t *testing.T) {
	mockUser := testmodels.NewLoginUser()
	lang := "en"

	mockResponse := mockUser.ToAdminResponse(lang)
	expected := models.AdminUserResponse{
		ID: mockUser.ID,
		Name: mockUser.Name,
		Username: mockUser.Username,
		Email: mockUser.Email,
		CountryCode: mockUser.Country.Code,
		Country: mockUser.Country.Names[lang],
		UserRole: mockUser.UserRole,
		IsActive: mockUser.IsActive,
	}

	assert.Equal(t, expected, mockResponse)
}

func TestToToken(t *testing.T) {
	mockUser := testmodels.NewLoginUser()

	mockToken := mockUser.ToToken()
	expected := models.UserToken{
		ID:       mockUser.ID,
		Username: mockUser.Username,
		UserRole: mockUser.UserRole,
	}

	assert.Equal(t, expected, mockToken)
}
