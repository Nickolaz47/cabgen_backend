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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAnalysisByID(t *testing.T) {
	testutils.SetupTestContext()

	mockAnalysis := testmodels.CreateMockAnalysis()
	mockResponse := mockAnalysis.ToResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockResponse,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: "abc1"}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.User.ID.String()}},
		)
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.User.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Analysis not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			FindByIDFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) (*models.AnalysisResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := analysis.NewAnalysisHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/analysis", "", nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.User.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.UserID})
		handler.GetAnalysisByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
