package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLaboratoryToAdminTableResponse(t *testing.T) {
	lab := models.Laboratory{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Abbreviation: "LACEN/RJ",
		IsActive:     true,
	}

	expected := models.LaboratoryAdminTableResponse{
		ID:           lab.ID,
		Name:         lab.Name,
		Abbreviation: lab.Abbreviation,
		IsActive:     lab.IsActive,
	}
	result := lab.ToAdminTableResponse()

	assert.Equal(t, expected, result)
}

func TestLaboratoryToFormResponse(t *testing.T) {
	lab := models.Laboratory{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Abbreviation: "LACEN/RJ",
		IsActive:     true,
	}

	expected := models.LaboratoryFormResponse{
		ID:           lab.ID,
		Name:         lab.Name,
		Abbreviation: lab.Abbreviation,
	}
	result := lab.ToFormResponse()

	assert.Equal(t, expected, result)
}
