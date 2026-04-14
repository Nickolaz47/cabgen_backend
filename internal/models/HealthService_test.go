package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHealthServiceTypeIsValid(t *testing.T) {
	tests := []struct {
		name     string
		hsType   models.HealthServiceType
		expected bool
	}{
		{
			name:     "Success - Public",
			hsType:   models.Public,
			expected: true,
		},
		{
			name:     "Success - Private",
			hsType:   models.Private,
			expected: true,
		},
		{
			name:     "Error",
			hsType:   "Hospital",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.hsType.IsValid())
		})
	}
}

func TestHealthServiceToAdminTableResponse(t *testing.T) {
	mockCountry := testmodels.NewCountry("", nil)

	healthService := models.HealthService{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Type:         "Public",
		CountryID:    mockCountry.ID,
		Country:      mockCountry,
		City:         "Rio de Janeiro",
		Contactant:   "John Doe",
		ContactEmail: "jonh@mail.com",
		ContactPhone: "123456789",
		IsActive:     true,
	}

	expected := models.HealthServiceAdminTableResponse{
		ID:           healthService.ID,
		Name:         healthService.Name,
		Type:         healthService.Type,
		Country:      healthService.Country.Code,
		City:         healthService.City,
		Contactant:   healthService.Contactant,
		ContactEmail: healthService.ContactEmail,
		ContactPhone: healthService.ContactPhone,
		IsActive:     healthService.IsActive,
	}

	result := healthService.ToAdminTableResponse()

	assert.Equal(t, expected, result)
}

func TestHealthServiceToFormResponse(t *testing.T) {
	mockCountry := testmodels.NewCountry("", nil)

	healthService := models.HealthService{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Type:         "Public",
		CountryID:    mockCountry.ID,
		Country:      mockCountry,
		City:         "Rio de Janeiro",
		Contactant:   "John Doe",
		ContactEmail: "jonh@mail.com",
		ContactPhone: "123456789",
		IsActive:     true,
	}

	expected := models.HealthServiceFormResponse{
		ID:   healthService.ID,
		Name: healthService.Name,
	}

	result := healthService.ToFormResponse()

	assert.Equal(t, expected, result)
}
