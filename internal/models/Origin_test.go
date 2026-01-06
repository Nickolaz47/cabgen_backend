package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOriginToResponse(t *testing.T) {
	mockOrigin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	lang := "en"

	result := mockOrigin.ToResponse(lang)
	expected := models.OriginResponse{
		Name:     mockOrigin.Names[lang],
		IsActive: mockOrigin.IsActive,
	}

	assert.Equal(t, expected, result)
}

func TestOriginToFormResponse(t *testing.T) {
	mockOrigin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	lang := "pt"

	result := mockOrigin.ToFormResponse(lang)
	expected := models.OriginFormResponse{
		ID:   mockOrigin.ID,
		Name: mockOrigin.Names[lang],
	}

	assert.Equal(t, expected, result)
}
