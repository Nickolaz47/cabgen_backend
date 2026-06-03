package analysis_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDownloadZip(t *testing.T) {
	testutils.SetupTestContext()

	mockUserID := uuid.New()
	rootDir := t.TempDir()

	zipPath := filepath.Join(rootDir, "analysis_results.zip")
	testutils.WriteMockFile(t, zipPath, []byte("test"))

	mockAnalysis := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.AnalysisAdminResponse, error) {
				response := mockAnalysis.ToAdminResponse()
				response.ResultsZipPath = &zipPath
				return &response, nil
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		router := gin.New()
		router.GET("/api/admin/analysis/:analysisId/download/zip",
			func(c *gin.Context) {
				c.Set("user", &models.UserToken{ID: mockUserID})
				handler.DownloadZip(c)
			})

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/admin/analysis/"+mockAnalysis.ID.String()+"/download/zip",
			nil,
		)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test", w.Body.String())
	})

	t.Run("Error - Zip Not Found", func(t *testing.T) {
		noZipAnalysis := testmodels.CreateMockAnalysis()
		noZipAnalysis.ResultsZipPath = nil

		svc := &mocks.MockAdminAnalysisService{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.AnalysisAdminResponse, error) {
				response := noZipAnalysis.ToAdminResponse()
				return &response, nil
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/analysis/:analysisId/download/zip",
			"",
			nil,
			nil,
		)
		c.Params = gin.Params{{Key: "analysisId",
			Value: noZipAnalysis.ID.String()}}
		handler.DownloadZip(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "The zip file is not available for this analysis.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/analysis/:analysisId/download/zip",
			"",
			nil,
			nil,
		)
		c.Params = gin.Params{{Key: "analysisId", Value: "not-a-valid-uuid"}}
		handler.DownloadZip(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "The URL ID is invalid.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.AnalysisAdminResponse, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/analysis/:analysisId/download/zip",
			"",
			nil,
			nil,
		)
		c.Params = gin.Params{{Key: "analysisId",
			Value: mockAnalysis.ID.String()}}
		handler.DownloadZip(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Analysis not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			FindByIDFunc: func(ctx context.Context, id uuid.UUID) (
				*models.AnalysisAdminResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/analysis/:analysisId/download/zip",
			"",
			nil,
			nil,
		)
		c.Params = gin.Params{{Key: "analysisId",
			Value: mockAnalysis.ID.String()}}
		handler.DownloadZip(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
