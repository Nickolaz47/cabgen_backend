package metrics_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/metrics"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetMetrics(t *testing.T) {
	testutils.SetupTestContext()

	mockResponse := &models.AdminMetricsResponse{
		PublicMetricsResponse: models.PublicMetricsResponse{
			TotalSamples:    1247,
			TotalCountries:  23,
			TotalSpecies:    156,
			TotalResistance: 342,
		},
	}

	expected := mockResponse.ToPublicResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMetricsService{
			GetMetricsFunc: func(ctx context.Context) (
				*models.AdminMetricsResponse, error) {
				return mockResponse, nil
			},
		}

		handler := metrics.NewMetricsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/metrics", "", nil, nil,
		)
		handler.GetMetrics(c)

		resp := testutils.ToJSON(map[string]any{"data": expected})

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

		handler := metrics.NewMetricsHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/metrics", "", nil, nil,
		)
		handler.GetMetrics(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
