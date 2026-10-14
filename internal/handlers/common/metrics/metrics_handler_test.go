package metrics_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/metrics"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserMetrics(t *testing.T) {
	testutils.SetupTestContext()

	mockUser := testmodels.NewLoginUser()

	t.Run("Success", func(t *testing.T) {
		expected := models.UserMetricsResponse{
			AnalysesByStatus: models.AnalysesByStatus{Done: 1},
			AnalysesByType: []models.AnalysisTypeMetric{
				{Type: string(models.AnalysisTypeFastQC), Count: 1},
			},
			TopSpecies: []models.SpeciesMetric{
				{Species: "Acinetobacter sp", Count: 1},
			},
			TopResistanceGenes: []models.GeneMetric{
				{Gene: "blaOXA-23", Count: 1},
			},
			SamplesByMonth: []models.MonthMetric{
				{Month: "2026-01", Count: 3},
			},
		}

		svc := &mocks.MockMetricsService{
			GetUserMetricsFunc: func(ctx context.Context,
				userID uuid.UUID) (*models.UserMetricsResponse, error) {
				assert.Equal(t, mockUser.ID, userID)
				return &expected, nil
			},
		}
		handler := metrics.NewMetricsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/users/me/metrics",
			"", nil, nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUser.ID})
		handler.GetUserMetrics(c)

		resp := testutils.ToJSON(
			map[string]any{"data": expected},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, resp, w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		handler := metrics.NewMetricsHandler(&mocks.MockMetricsService{})

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/users/me/metrics",
			"", nil, nil,
		)
		handler.GetUserMetrics(c)

		resp := testutils.ToJSON(
			map[string]string{
				"error": "Unauthorized. Please log in to continue.",
			},
		)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, resp, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockMetricsService{
			GetUserMetricsFunc: func(ctx context.Context,
				userID uuid.UUID) (*models.UserMetricsResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := metrics.NewMetricsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/users/me/metrics",
			"", nil, nil,
		)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUser.ID})
		handler.GetUserMetrics(c)

		resp := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, resp, w.Body.String())
	})
}
