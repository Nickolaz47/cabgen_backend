package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSampleSourceToAdminDetailResponse(t *testing.T) {
	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)

	expected := models.SampleSourceAdminDetailResponse{
		ID:       mockSampleSource.ID,
		Names:    mockSampleSource.Names,
		Groups:   mockSampleSource.Groups,
		IsActive: mockSampleSource.IsActive,
	}
	result := mockSampleSource.ToAdminDetailResponse()

	assert.Equal(t, expected, result)
}

func TestSampleSourceToAdminTableResponse(t *testing.T) {
	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)

	language := "en"
	expected := models.SampleSourceAdminTableResponse{
		ID:       mockSampleSource.ID,
		Name:     mockSampleSource.Names[language],
		Group:    mockSampleSource.Groups[language],
		IsActive: mockSampleSource.IsActive,
	}
	result := mockSampleSource.ToAdminTableResponse(language)

	assert.Equal(t, expected, result)
}

func SampleSourceToFormResponse(t *testing.T) {
	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)

	expected := models.SampleSourceFormResponse{
		ID:   mockSampleSource.ID,
		Name: mockSampleSource.Names["en"],
	}
	result := mockSampleSource.ToFormResponse("en")

	assert.Equal(t, expected, result)
}
