package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func AnalysisToAdminResponse(t *testing.T) {
	mockAnalysis := testmodels.CreateMockAnalysis()

	expected := models.AnalysisAdminResponse{
		ID:           mockAnalysis.ID,
		Type:         mockAnalysis.Type,
		Status:       mockAnalysis.Status,
		ErrorMessage: mockAnalysis.ErrorMessage,
		Sample:       mockAnalysis.Sample.Name,
		SampleID:     mockAnalysis.Sample.ID,
		User:         mockAnalysis.User.Username,
		UserID:       mockAnalysis.UserID,
		Metrics:      mockAnalysis.Metrics,
		FastQC1:      mockAnalysis.FastQC1,
		FastQC2:      mockAnalysis.FastQC2,
		StartedAt:    mockAnalysis.StartedAt,
		FinishedAt:   mockAnalysis.FinishedAt,
	}
	result := mockAnalysis.ToAdminResponse()

	assert.Equal(t, expected, result)
}

func AnalysisToResponse(t *testing.T) {
	mockAnalysis := testmodels.CreateMockAnalysis()

	expected := models.AnalysisResponse{
		ID:           mockAnalysis.ID,
		Type:         mockAnalysis.Type,
		Status:       mockAnalysis.Status,
		ErrorMessage: mockAnalysis.ErrorMessage,
		Sample:       mockAnalysis.Sample.Name,
		SampleID:     mockAnalysis.Sample.ID,
		Metrics:      mockAnalysis.Metrics,
		FastQC1:      mockAnalysis.FastQC1,
		FastQC2:      mockAnalysis.FastQC2,
		StartedAt:    mockAnalysis.StartedAt,
		FinishedAt:   mockAnalysis.FinishedAt,
	}
	result := mockAnalysis.ToResponse()

	assert.Equal(t, expected, result)
}
