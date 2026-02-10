package models_test

import (
	"fmt"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMicroorganismToAdminDetailResponse(t *testing.T) {
	id := uuid.NewString()
	taxon := models.Bacteria
	species := "Neisseria meningitidis"
	varietyMap := map[string]string{
		"pt": "Sorogrupo B",
		"en": "Serogroup B",
		"es": "Serogrupo B",
	}

	mockMicro := testmodels.NewMicroorganism(
		id, taxon, species, varietyMap, true,
	)

	expected := models.MicroorganismAdminDetailResponse{
		ID:       mockMicro.ID,
		Taxon:    mockMicro.Taxon,
		Species:  mockMicro.Species,
		Variety:  mockMicro.Variety,
		IsActive: mockMicro.IsActive,
	}
	result := mockMicro.ToAdminDetailResponse()

	assert.Equal(t, expected, result)
}

func TestMicroorganismToAdminTableResponse(t *testing.T) {
	id := uuid.NewString()
	taxon := models.Bacteria
	species := "Neisseria meningitidis"
	varietyMap := map[string]string{
		"pt": "Sorogrupo B",
		"en": "Serogroup B",
		"es": "Serogrupo B",
	}

	mockMicro := testmodels.NewMicroorganism(
		id, taxon, species, varietyMap, true,
	)

	lang := "pt"

	expected := models.MicroorganismAdminTableResponse{
		ID:       mockMicro.ID,
		Taxon:    mockMicro.Taxon,
		Species:  mockMicro.Species,
		Variety:  mockMicro.Variety[lang],
		IsActive: mockMicro.IsActive,
	}
	result := mockMicro.ToAdminTableResponse(lang)

	assert.Equal(t, expected, result)
}

func TestMicroorganismToFormResponse(t *testing.T) {
	id := uuid.NewString()
	taxon := models.Bacteria
	species := "Neisseria meningitidis"
	varietyMap := map[string]string{
		"pt": "Sorogrupo B",
		"en": "Serogroup B",
		"es": "Serogrupo B",
	}

	mockMicro := testmodels.NewMicroorganism(
		id, taxon, species, varietyMap, true,
	)

	lang := "pt"
	expectedName := fmt.Sprintf("%s %s", mockMicro.Species,
		mockMicro.Variety[lang])

	expected := models.MicroorganismFormResponse{
		ID:      mockMicro.ID,
		Species: expectedName,
	}
	result := mockMicro.ToFormResponse(lang)

	assert.Equal(t, expected, result)
}
