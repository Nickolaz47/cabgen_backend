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

func TestOriginToResponse(t *testing.T) {
	mockOrigin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

	result := mockOrigin.ToResponse(c)
	expected := models.OriginResponse{
		ID:       mockOrigin.ID,
		Name:     mockOrigin.Names["en"],
		IsActive: mockOrigin.IsActive,
	}

	assert.Equal(t, expected, result)
}

func TestOriginToFormResponse(t *testing.T) {
	mockOrigin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	c, _ := testutils.SetupGinContext(http.MethodGet, "/", "", nil, nil)

	result := mockOrigin.ToFormResponse(c)
	expected := models.OriginFormResponse{
		ID:   mockOrigin.ID,
		Name: mockOrigin.Names["en"],
	}

	assert.Equal(t, expected, result)
}
