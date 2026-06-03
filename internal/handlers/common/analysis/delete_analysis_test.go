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

func TestDeleteAnalysis(t *testing.T) {
	testutils.SetupTestContext()

	mockAnalysis := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}

		handler := analysis.NewAnalysisHandler(svc)
		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/analysis",
			"",
			nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.User.ID})
		handler.DeleteAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Analysis deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{}

		handler := analysis.NewAnalysisHandler(svc)
		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/analysis",
			"",
			nil,
			gin.Params{{Key: "analysisId", Value: "123"}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.User.ID})
		handler.DeleteAnalysis(c)

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
			http.MethodDelete,
			"/api/analysis",
			"",
			nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)

		handler.DeleteAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockAnalysisService{
			DeleteFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) error {
				return services.ErrNotFound
			},
		}

		handler := analysis.NewAnalysisHandler(svc)
		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/analysis",
			"",
			nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.User.ID})

		handler.DeleteAnalysis(c)

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
			DeleteFunc: func(ctx context.Context, analysisID,
				userID uuid.UUID) error {
				return services.ErrInternal
			},
		}

		handler := analysis.NewAnalysisHandler(svc)
		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/analysis",
			"",
			nil,
			gin.Params{{Key: "analysisId", Value: mockAnalysis.ID.String()}},
		)
		c.Set("user", &models.UserToken{ID: mockAnalysis.User.ID})

		handler.DeleteAnalysis(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
