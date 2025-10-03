package models_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestUserToResponse(t *testing.T) {
	mockUser := testmodels.NewLoginUser()
	c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

	mockResponse := mockUser.ToResponse(c)
	expected := models.UserResponse{
		ID:          mockUser.ID,
		Name:        mockUser.Name,
		Username:    mockUser.Username,
		Email:       mockUser.Email,
		CountryCode: mockUser.CountryCode,
		Country:     mockUser.Country.Names["en"],
		UserRole:    mockUser.UserRole,
		Interest:    mockUser.Interest,
		Role:        mockUser.Role,
		Institution: mockUser.Institution,
	}

	assert.Equal(t, expected, mockResponse)
}

func TestToToken(t *testing.T) {
	mockUser := testmodels.NewLoginUser()

	mockToken := mockUser.ToToken()
	expected := models.UserToken{
		ID: mockUser.ID,
		Username: mockUser.Username,
		UserRole: mockUser.UserRole,
	}

	assert.Equal(t, expected, mockToken)
}
