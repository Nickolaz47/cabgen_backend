package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
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
	result := mockAnalysis.ToResponse("en")

	assert.Equal(t, expected, result)
}

func TestAnalysisToResponseTranslation(t *testing.T) {
	t.Run("Translates known error to PT", func(t *testing.T) {
		msg := pipeline.ErrFastQC.Error()
		mock := testmodels.CreateMockAnalysis()
		mock.ErrorMessage = &msg

		result := mock.ToResponse("pt")
		assert.NotNil(t, result.ErrorMessage)
		assert.Equal(t, "A etapa do FastQC falhou. Crie uma nova análise.",
			*result.ErrorMessage)
	})

	t.Run("Translates known error to ES", func(t *testing.T) {
		msg := pipeline.ErrCorruptedInput.Error()
		mock := testmodels.CreateMockAnalysis()
		mock.ErrorMessage = &msg

		result := mock.ToResponse("es")
		assert.NotNil(t, result.ErrorMessage)
		assert.Equal(t, "El archivo de entrada está corrupto o truncado.",
			*result.ErrorMessage)
	})

	t.Run("Unknown error stays as-is", func(t *testing.T) {
		msg := "some unexpected error"
		mock := testmodels.CreateMockAnalysis()
		mock.ErrorMessage = &msg

		result := mock.ToResponse("pt")
		assert.NotNil(t, result.ErrorMessage)
		assert.Equal(t, "some unexpected error", *result.ErrorMessage)
	})

	t.Run("Nil ErrorMessage stays nil", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.ErrorMessage = nil

		result := mock.ToResponse("en")
		assert.Nil(t, result.ErrorMessage)
	})
}
