package country_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/country"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCountry(t *testing.T) {
	testutils.SetupTestContext()

	code := "BRA"
	names := map[string]string{"pt": "Brasil", "en": "Brazil", "es": "Brazil"}

	updateInput := models.CountryUpdateInput{
		Code:  &code,
		Names: names,
	}

	mockCountry := testmodels.NewCountry(code, names)
	mockResponse := mockCountry.ToAdminDetailResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			UpdateFunc: func(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error) {
				return &mockResponse, nil
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/country",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "code", Value: "BRA"}},
		)
		handler.UpdateCountry(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockResponse,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateCountryTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockCountryService{}
			handler := country.NewAdminCountryHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPut,
				"/api/admin/country",
				tt.Body,
				nil,
				gin.Params{{Key: "code", Value: "BRA"}},
			)
			handler.UpdateCountry(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			UpdateFunc: func(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/country",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "code", Value: "BRA"}},
		)
		handler.UpdateCountry(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "No country was found with this code.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			UpdateFunc: func(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/country",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "code", Value: "BRA"}},
		)
		handler.UpdateCountry(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Country already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			UpdateFunc: func(ctx context.Context, code string, input models.CountryUpdateInput) (*models.CountryAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/country",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "code", Value: "BRA"}},
		)
		handler.UpdateCountry(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
