package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOriginToAdminDetailResponse(t *testing.T) {
	mockOrigin := testmodels.NewOrigin(uuid.New().String(),
		map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		true,
	)

	result := mockOrigin.ToAdminDetailResponse()
	expected := models.OriginAdminDetailResponse{
		ID:       mockOrigin.ID,
		Names:    mockOrigin.Names,
		IsActive: mockOrigin.IsActive,
	}

	assert.Equal(t, expected, result)
}

func TestOriginToAdminTableResponse(t *testing.T) {
	mockOrigin := testmodels.NewOrigin(uuid.New().String(),
		map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		true,
	)
	lang := "en"

	result := mockOrigin.ToAdminTableResponse(lang)
	expected := models.OriginAdminTableResponse{
		ID:       mockOrigin.ID,
		Name:     mockOrigin.Names[lang],
		IsActive: mockOrigin.IsActive,
	}

	assert.Equal(t, expected, result)
}

func TestOriginToFormResponse(t *testing.T) {
	mockOrigin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		true)
	lang := "pt"

	result := mockOrigin.ToFormResponse(lang)
	expected := models.OriginFormResponse{
		ID:   mockOrigin.ID,
		Name: mockOrigin.Names[lang],
	}

	assert.Equal(t, expected, result)
}
