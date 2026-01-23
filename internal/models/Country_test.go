package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestCountryToAdminDetailResponse(t *testing.T) {
	mockCountry := testmodels.NewCountry("", nil)

	expected := models.CountryAdminDetailResponse{
		Code:  mockCountry.Code,
		Names: mockCountry.Names,
	}
	result := mockCountry.ToAdminDetailResponse()

	assert.Equal(t, expected, result)
}

func TestCountryToFormDetailResponse(t *testing.T) {
	mockCountry := testmodels.NewCountry("", nil)
	lang := "en"

	expected := models.CountryFormResponse{
		Code: mockCountry.Code,
		Name: mockCountry.Names[lang],
	}
	result := mockCountry.ToFormResponse(lang)

	assert.Equal(t, expected, result)
}
