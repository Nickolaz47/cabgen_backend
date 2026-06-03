package analysis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAnalysis(t *testing.T) {
	testutils.SetupTestContext()

	const validUUID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	mockAnalysis := testmodels.CreateMockAnalysis()
	mockResponse := mockAnalysis.ToAdminResponse()

	validInput := map[string]any{
		"status":           models.AnalysisStatusDone,
		"metrics":          map[string]any{"coverage": 98.5, "reads": 1500000},
		"fastqc1":          "/app/uploads/fastqc1_report.html",
		"fastqc2":          "/app/uploads/fastqc2_report.html",
		"results_zip_path": "/app/results/analysis_results.zip",
		"error_message":    nil,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			UpdateFunc: func(ctx context.Context, analysisID uuid.UUID,
				input models.AdminAnalysisUpdateInput) (
				*models.AnalysisAdminResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/analysis",
			testutils.ToJSON(validInput),
			nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		handler.UpdateAnalysis(c)

		expected := testutils.ToJSON(
			map[string]any{"data": mockResponse},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Type", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{}
		handler := analysis.NewAdminAnalysisHandler(svc)

		invalidInput := testutils.CopyMap(validInput)
		invalidInput["status"] = "invalid"

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/analysis",
			testutils.ToJSON(invalidInput),
			nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		handler.UpdateAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "This analysis status is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{}
		handler := analysis.NewAdminAnalysisHandler(svc)

		for _, test := range data.AdminAnalysisUpdateTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPost,
					"/api/admin/analysis",
					test.Body,
					nil,
					gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
				)
				handler.UpdateAnalysis(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			UpdateFunc: func(ctx context.Context, analysisID uuid.UUID,
				input models.AdminAnalysisUpdateInput) (
				*models.AnalysisAdminResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/analysis",
			testutils.ToJSON(validInput),
			nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		handler.UpdateAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
