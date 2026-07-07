package city_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/city"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetCities(t *testing.T) {
	testutils.SetupTestContext()

	mockResponse := []models.SelectOption{
		{Label: "Belo Horizonte - MG", Value: "Belo Horizonte - MG"},
		{Label: "Maricá - RJ", Value: "Maricá - RJ"},
		{Label: "São Paulo - SP", Value: "São Paulo - SP"},
		{Label: "option.city.other", Value: "Other"},
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockCityService{
			FindAllFunc: func(ctx context.Context) (
				[]models.SelectOption, error) {
				return mockResponse, nil
			},
		}
		handler := city.NewCityHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/cities",
			"", nil, nil,
		)
		handler.GetCities(c)

		expected := testutils.ToJSON(map[string]any{
			"data": mockResponse,
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockCityService{
			FindAllFunc: func(ctx context.Context) ([]models.SelectOption, error) {
				return nil, services.ErrInternal
			},
		}
		handler := city.NewCityHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/cities",
			"", nil, nil,
		)
		handler.GetCities(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
