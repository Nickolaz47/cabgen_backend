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
	"github.com/stretchr/testify/assert"
)

func TestCreateCountry(t *testing.T) {
	testutils.SetupTestContext()

	code := "BRA"
	names := map[string]string{"pt": "Brasil", "en": "Brazil", "es": "Brazil"}

	createInput := models.CountryCreateInput{
		Code:  code,
		Names: names,
	}

	mockCountry := testmodels.NewCountry(code, names)

	mockResponse := mockCountry.ToAdminDetailResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			CreateFunc: func(ctx context.Context, input models.CountryCreateInput) (*models.CountryAdminDetailResponse, error) {
				return &mockResponse, nil
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		body := testutils.ToJSON(createInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/country", body,
			nil, nil,
		)
		handler.CreateCountry(c)

		expected := testutils.ToJSON(
			map[string]any{
				"message": "Country created successfully.",
				"data":    mockResponse,
			},
		)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.CreateCountryTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockCountryService{}
			handler := country.NewAdminCountryHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/country", tt.Body, nil, nil,
			)
			handler.CreateCountry(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			CreateFunc: func(ctx context.Context, input models.CountryCreateInput) (*models.CountryAdminDetailResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		body := testutils.ToJSON(createInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/country", body,
			nil, nil,
		)
		handler.CreateCountry(c)

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
			CreateFunc: func(ctx context.Context, input models.CountryCreateInput) (*models.CountryAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		body := testutils.ToJSON(createInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/country", body,
			nil, nil,
		)
		handler.CreateCountry(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
