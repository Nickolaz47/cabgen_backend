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

func TestTaxonToTranslatedString(t *testing.T) {
	tests := []struct {
		name     string
		language string
		taxon    models.Taxon
		expected string
	}{
		{
			name:     "Bacteria to portuguese",
			language: "pt",
			taxon:    models.Bacteria,
			expected: "Bactéria",
		},
		{
			name:     "Fungi to spanish",
			language: "es",
			taxon:    models.Fungi,
			expected: "Hongos",
		},
		{
			name:     "Virus to portuguese",
			language: "pt",
			taxon:    models.Virus,
			expected: "Vírus",
		},
		{
			name:     "Fungi to english",
			language: "en",
			taxon:    models.Fungi,
			expected: "Fungi",
		},
		{
			name:     "Invalid language",
			language: "an",
			taxon:    models.Bacteria,
			expected: "Bacteria",
		},
		{
			name:     "Unknown taxon",
			language: "pt",
			taxon:    models.Taxon("Ghost"),
			expected: "Ghost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.taxon.ToTranslatedString(tt.language)

			assert.Equal(t, tt.expected, result)
		})
	}
}
