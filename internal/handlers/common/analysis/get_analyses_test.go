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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAnalyses(t *testing.T) {
	testutils.SetupTestContext()

	mockAnalysis := testmodels.CreateMockAnalysis()
	mockResponse := mockAnalysis.ToResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindAllFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.AnalysisResponse, error) {
				return []models.AnalysisResponse{mockResponse}, nil
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil, nil,
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalyses(c)

		expected := testutils.ToJSON(
			map[string][]models.AnalysisResponse{
				"data": {mockResponse},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil, nil,
		)
		handler.GetAnalyses(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindAllFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.AnalysisResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil, nil,
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
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
