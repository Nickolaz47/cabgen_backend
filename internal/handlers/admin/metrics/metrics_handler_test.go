package metrics_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/metrics"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAdminGetMetrics(t *testing.T) {
	testutils.SetupTestContext()

	mockResponse := &models.AdminMetricsResponse{
		PublicMetricsResponse: models.PublicMetricsResponse{
			TotalSamples:    1247,
			TotalCountries:  23,
			TotalSpecies:    156,
			TotalResistance: 342,
		},
		TotalUsers: 89,
		AnalysesByStatus: models.AnalysesByStatus{
			Done: 2890, Running: 12, Pending: 456, Failed: 54,
		},
		TopCountries: []models.CountryMetric{
			{Country: "BRA", Count: 412},
		},
		SpeciesBreakdown: []models.SpeciesMetric{
			{Species: "Acinetobacter baumannii", Count: 312},
		},
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMetricsService{
			GetMetricsFunc: func(ctx context.Context) (
				*models.AdminMetricsResponse, error) {
				return mockResponse, nil
			},
		}

		handler := metrics.NewAdminMetricsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/metrics", "", nil, nil,
		)
		handler.GetMetrics(c)

		resp := testutils.ToJSON(map[string]any{"data": mockResponse})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, resp, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockMetricsService{
			GetMetricsFunc: func(ctx context.Context) (
				*models.AdminMetricsResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := metrics.NewAdminMetricsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/metrics", "", nil, nil,
		)
		handler.GetMetrics(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
