package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToFormResponse(t *testing.T) {
	sequencer := models.Sequencer{
		ID:       uuid.New(),
		Brand:    "Illumina",
		Model:    "MySeq",
		IsActive: true,
	}

	expected := models.SequencerFormResponse{
		ID:    sequencer.ID,
		Brand: sequencer.Brand,
	}
	result := sequencer.ToFormResponse()

	assert.Equal(t, expected, result)
}
