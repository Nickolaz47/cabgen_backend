package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestAnalysisToResponse(t *testing.T) {
	mockAnalysis := testmodels.CreateMockAnalysis()
	mockAnalysis.Step = models.StepCheckM

	expected := models.AnalysisResponse{
		ID:             mockAnalysis.ID,
		Type:           mockAnalysis.Type,
		Status:         mockAnalysis.Status,
		Step:           mockAnalysis.Step,
		ErrorMessage:   mockAnalysis.ErrorMessage,
		Sample:         mockAnalysis.Sample.OriginCode,
		SampleID:       mockAnalysis.Sample.ID,
		User:           mockAnalysis.User.Username,
		UserID:         mockAnalysis.UserID,
		Metrics:        mockAnalysis.Metrics,
		ResultsZipPath: mockAnalysis.ResultsZipPath,
		FastQC1:        mockAnalysis.FastQC1,
		FastQC2:        mockAnalysis.FastQC2,
		StartedAt:      mockAnalysis.StartedAt,
		FinishedAt:     mockAnalysis.FinishedAt,
	}
	result := mockAnalysis.ToResponse()

	assert.Equal(t, expected, result)
}
