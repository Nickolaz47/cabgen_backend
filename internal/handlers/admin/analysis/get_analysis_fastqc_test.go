package analysis_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminGetAnalysisFastQCByID(t *testing.T) {
	testutils.SetupTestContext()

	analysisID := uuid.New()
	reportContent := "fastqc html report"
	reportFile := filepath.Join(t.TempDir(), "fastqc.html")
	os.WriteFile(reportFile, []byte(reportContent), 0644)

	t.Run("Success - FastQC1", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID, userID uuid.UUID,
				language string) (*models.AnalysisResponse, error) {
				return &models.AnalysisResponse{
					ID:      analysisID,
					Sample:  "A01",
					Type:    models.AnalysisTypeFastQC,
					FastQC1: &reportFile,
				}, nil
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/admin/analyses",
			"", nil, gin.Params{
				{Key: "analysisId", Value: analysisID.String()},
				{Key: "fastqcReport", Value: "fastqc1"},
			})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, reportContent, w.Body.String())
	})

	t.Run("Success - FastQC2", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID, userID uuid.UUID,
				language string) (*models.AnalysisResponse, error) {
				return &models.AnalysisResponse{
					ID:      analysisID,
					FastQC2: &reportFile,
				}, nil
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/admin/analyses",
			"", nil, gin.Params{
				{Key: "analysisId", Value: analysisID.String()},
				{Key: "fastqcReport", Value: "fastqc2"},
			})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, reportContent, w.Body.String())
	})

	t.Run("Error - Invalid fastqcReport", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID, userID uuid.UUID,
				language string) (*models.AnalysisResponse, error) {
				return &models.AnalysisResponse{ID: analysisID}, nil
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/admin/analyses",
			"", nil, gin.Params{
				{Key: "analysisId", Value: analysisID.String()},
				{Key: "fastqcReport", Value: "fastqc3"},
			})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid FastQC report type.")
	})

	t.Run("Error - Invalid Analysis ID", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/admin/analyses",
			"", nil, gin.Params{
				{Key: "analysisId", Value: "not-a-uuid"},
				{Key: "fastqcReport", Value: "fastqc1"},
			})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "The URL ID is invalid.")
	})

	t.Run("Error - Analysis Not Found", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID, userID uuid.UUID,
				language string) (*models.AnalysisResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/admin/analyses",
			"", nil, gin.Params{
				{Key: "analysisId", Value: uuid.NewString()},
				{Key: "fastqcReport", Value: "fastqc1"},
			})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Analysis not found.")
	})

	t.Run("Error - Report Not Available", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID, userID uuid.UUID,
				language string) (*models.AnalysisResponse, error) {
				return &models.AnalysisResponse{ID: analysisID}, nil
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/admin/analyses",
			"", nil, gin.Params{
				{Key: "analysisId", Value: analysisID.String()},
				{Key: "fastqcReport", Value: "fastqc1"},
			})
		handler.GetAnalysisFastQCByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "FastQC report")
	})
}
