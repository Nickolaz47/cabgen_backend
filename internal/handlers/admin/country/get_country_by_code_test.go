package country_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/country"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCountryByCode(t *testing.T) {
	testutils.SetupTestContext()

	mockCountry := testmodels.NewCountry("", nil)
	mockCountry.ID = 1

	response := mockCountry.ToAdminDetailResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			FindByCodeFunc: func(ctx context.Context, code string) (*models.CountryAdminDetailResponse, error) {
				return &response, nil
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/country",
			"", nil, gin.Params{{Key: "code", Value: mockCountry.Code}},
		)
		handler.GetCountryByCode(c)

		expected := testutils.ToJSON(
			map[string]models.CountryAdminDetailResponse{
				"data": response,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			FindByCodeFunc: func(ctx context.Context, code string) (*models.CountryAdminDetailResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/country",
			"", nil, gin.Params{{Key: "code", Value: "ARG"}},
		)
		handler.GetCountryByCode(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "No country was found with this code.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			FindByCodeFunc: func(ctx context.Context, code string) (*models.CountryAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/country",
			"", nil, gin.Params{{Key: "code", Value: mockCountry.Code}},
		)
		handler.GetCountryByCode(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
