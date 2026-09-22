package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
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
		ID:          mockUser.ID,
		Name:        mockUser.Name,
		Username:    mockUser.Username,
		Email:       mockUser.Email,
		CountryCode: mockUser.Country.Code,
		Country:     mockUser.Country.Names[lang],
		UserRole:    mockUser.UserRole,
		IsActive:    mockUser.IsActive,
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

func TestAdminUserFilterBinding(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected models.AdminUserFilter
		wantErr  bool
	}{
		{
			name:  "Success - Full filter",
			query: "input=jo&userRole=Admin&active=true",
			expected: models.AdminUserFilter{
				Input:    "jo",
				UserRole: models.Admin,
				Active: func() *bool {
					active := true
					return &active
				}(),
			},
		},
		{
			name:     "Success - Empty query",
			query:    "",
			expected: models.AdminUserFilter{},
		},
		{
			name:    "Error - Invalid active flag",
			query:   "active=yes",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := testutils.BindFilter[models.AdminUserFilter](tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tt.expected, filter)
			}
		})
	}
}
