package analysis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestDownloadDashboardTSV(t *testing.T) {
	testutils.SetupTestContext()

	mockAnalysis := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			DownloadDashboardTSVFunc: func(ctx context.Context) (
				[]models.Analysis, error) {
				return []models.Analysis{mockAnalysis}, nil
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodPost,
			"/api/admin/analyses/download/dashboard", "", nil, nil)
		handler.DownloadDashboardTSV(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "text/tab-separated-values; charset=utf-8",
			w.Header().Get("Content-Type"))
		assert.Equal(t, "attachment; filename=cabgen_dashboard.tsv",
			w.Header().Get("Content-Disposition"))

		body := w.Body.String()
		assert.Contains(t, body,
			"ID\tEspécie\tPlasmídeos\tMLST\tCarbapenemases")
		assert.Contains(t, body, mockAnalysis.Sample.ID.String())
		assert.Contains(t, body, "Acinetobacter sp")
		assert.Contains(t, body, "ST502")
		assert.Contains(t, body, "Aspirated")
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			DownloadDashboardTSVFunc: func(ctx context.Context) (
				[]models.Analysis, error) {
				return nil, services.ErrInternal
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodPost,
			"/api/admin/analyses/download/dashboard", "", nil, nil)
		handler.DownloadDashboardTSV(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - No Dashboard Data", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			DownloadDashboardTSVFunc: func(ctx context.Context) (
				[]models.Analysis, error) {
				return nil, services.ErrNoDashboardData
			},
		}
		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodPost,
			"/api/admin/analyses/download/dashboard", "", nil, nil)
		handler.DownloadDashboardTSV(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "No analyses from in-network samples available for download.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
