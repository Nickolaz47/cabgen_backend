package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSequencerToFormResponse(t *testing.T) {
	sequencer := models.Sequencer{
		ID:       uuid.New(),
		Brand:    "Illumina",
		Model:    "MySeq",
		IsActive: true,
	}

	expected := models.SequencerFormResponse{
		ID:    sequencer.ID,
		Model: sequencer.Model,
		Brand: sequencer.Brand,
	}
	result := sequencer.ToFormResponse()

	assert.Equal(t, expected, result)
}

func TestSequencerToAdminTableResponse(t *testing.T) {
	sequencer := models.Sequencer{
		ID:       uuid.New(),
		Brand:    "Illumina",
		Model:    "MySeq",
		IsActive: true,
	}

	expected := models.SequencerAdminTableResponse{
		ID:       sequencer.ID,
		Model:    sequencer.Model,
		Brand:    sequencer.Brand,
		IsActive: sequencer.IsActive,
	}
	result := sequencer.ToAdminTableResponse()

	assert.Equal(t, expected, result)
}
