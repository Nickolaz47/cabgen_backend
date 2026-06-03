package analysis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDownloadBatchTSV(t *testing.T) {
	testutils.SetupTestContext()

	mockUserID := uuid.New()
	mockAnalysis := testmodels.CreateMockAnalysis()
	mockAnalyses := []models.AnalysisResponse{
		mockAnalysis.ToResponse(),
		mockAnalysis.ToResponse(),
	}

	const validUUID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	validInput := map[string]any{
		"ids": []string{validUUID, validUUID},
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindManyByIDsFunc: func(ctx context.Context, ids []uuid.UUID,
				userID uuid.UUID) ([]models.AnalysisResponse, error) {
				assert.Equal(t, mockUserID, userID)
				assert.Len(t, ids, 2)
				return mockAnalyses, nil
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis/download/tsv",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

		handler.DownloadBatchTSV(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "text/tab-separated-values",
			w.Header().Get("Content-Type"))
		assert.Equal(t, "attachment; filename=cabgen_results.tsv",
			w.Header().Get("Content-Disposition"))
		body := w.Body.String()

		assert.Contains(t, body, "id")
		assert.Contains(t, body, "type")
		assert.Contains(t, body, "status")
		assert.Contains(t, body, "started_at")

		mockResponse := mockAnalysis.ToResponse()
		assert.Contains(t, body, mockResponse.ID.String())
		assert.Contains(t, body, string(models.AnalysisTypeComplete))
		assert.Contains(t, body, string(models.AnalysisStatusDone))
		assert.Contains(t, body, "/result/fastqc_reads1.html")
		assert.Contains(t, body, "/result/fastqc_reads2.html")
		assert.Contains(t, body, "11-05-2024 00:00:00")
		assert.Contains(t, body, "95.89")
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		for _, test := range data.AnalysisTSVDownloadTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPost,
					"/api/analysis/download/tsv",
					test.Body,
					nil,
					nil,
				)
				c.Set("user", &models.UserToken{ID: mockUserID})

				handler.DownloadBatchTSV(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis/download/tsv",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)

		handler.DownloadBatchTSV(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Unauthorized. Please log in to continue.",
		})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindManyByIDsFunc: func(ctx context.Context, ids []uuid.UUID,
				userID uuid.UUID) ([]models.AnalysisResponse, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis/download/tsv",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

		handler.DownloadBatchTSV(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Analysis not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindManyByIDsFunc: func(ctx context.Context, ids []uuid.UUID,
				userID uuid.UUID) ([]models.AnalysisResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/analysis/download/tsv",
			testutils.ToJSON(validInput),
			nil,
			nil,
		)
		c.Set("user", &models.UserToken{ID: mockUserID})

		handler.DownloadBatchTSV(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
