package models_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SampleSourceToResponse(t *testing.T) {
	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)
	c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

	expected := models.SampleSourceResponse{
		Name:     mockSampleSource.Names["en"],
		Group:    mockSampleSource.Groups["en"],
		IsActive: mockSampleSource.IsActive,
	}
	result := mockSampleSource.ToResponse(c)

	assert.Equal(t, expected, result)
}

func SampleSourceToFormResponse(t *testing.T) {
	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)
	c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

	expected := models.SampleSourceFormResponse{
		ID:   mockSampleSource.ID,
		Name: mockSampleSource.Names["en"],
	}
	result := mockSampleSource.ToFormResponse(c)

	assert.Equal(t, expected, result)
}
