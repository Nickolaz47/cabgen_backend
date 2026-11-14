package analysis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAnalyses(t *testing.T) {
	testutils.SetupTestContext()

	mockUserID := uuid.New()
	mockSampleID := uuid.New()
	mockAnalysis := testmodels.CreateMockAnalysis()
	mockResponse := mockAnalysis.ToResponse("en")

	validInput := map[string]any{
		"type":       models.AnalysisTypeComplete,
		"sample_ids": []uuid.UUID{mockSampleID},
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			CreateBatchFunc: func(ctx context.Context,
				input models.AnalysisCreateBatchDTO,
				language string) ([]models.AnalysisResponse, error) {
				assert.Equal(t, mockUserID, input.UserID)
				assert.Equal(t, []uuid.UUID{mockSampleID}, input.SampleIDs)
				return []models.AnalysisResponse{mockResponse}, nil
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analyses/batch",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateAnalyses(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data":    []models.AnalysisResponse{mockResponse},
				"message": "Analyses created successfully.",
			},
		)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Type", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		invalidInput := testutils.CopyMap(validInput)
		invalidInput["type"] = "invalid_type"

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analyses/batch",
			testutils.ToJSON(invalidInput),
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "This analysis type is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		invalidInput := testutils.CopyMap(validInput)
		delete(invalidInput, "sample_ids")

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analyses/batch",
			testutils.ToJSON(invalidInput),
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The sample IDs are required.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analyses/batch",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)

		handler.CreateAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Exceeded Batch", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			CreateBatchFunc: func(ctx context.Context,
				input models.AnalysisCreateBatchDTO,
				language string) ([]models.AnalysisResponse, error) {
				return nil, services.ErrExceededAnalysisBatch
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analyses/batch",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The batch exceeds the maximum number of analyses.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Duplicate Sample", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			CreateBatchFunc: func(ctx context.Context,
				input models.AnalysisCreateBatchDTO,
				language string) ([]models.AnalysisResponse, error) {
				return nil, &services.AnalysisBatchSampleError{
					SampleID: mockSampleID,
					Err:      services.ErrDuplicateSample,
				}
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analyses/batch",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Duplicate sample ID in the batch: " +
					mockSampleID.String() + ".",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			CreateBatchFunc: func(ctx context.Context,
				input models.AnalysisCreateBatchDTO,
				language string) ([]models.AnalysisResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analyses/batch",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
