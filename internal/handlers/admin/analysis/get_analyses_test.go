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

func TestGetAnalyses(t *testing.T) {
	testutils.SetupTestContext()

	mockAnalysis := testmodels.CreateMockAnalysis()
	mockResponse := mockAnalysis.ToAdminResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			FindAllFunc: func(ctx context.Context) (
				[]models.AnalysisAdminResponse, error) {
				return []models.AnalysisAdminResponse{mockResponse}, nil
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/analysis", "", nil, nil,
		)
		handler.GetAnalyses(c)

		expected := testutils.ToJSON(
			map[string][]models.AnalysisAdminResponse{
				"data": {mockResponse},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAdminAnalysisService{
			FindAllFunc: func(ctx context.Context) (
				[]models.AnalysisAdminResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := analysis.NewAdminAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/analysis", "", nil, nil,
		)
		handler.GetAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

}
