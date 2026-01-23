package country_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/country"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetCountriesByName(t *testing.T) {
	testutils.SetupTestContext()

	mockCountry := testmodels.NewCountry("", nil)
	mockCountry.ID = 1

	response := mockCountry.ToFormResponse("pt")

	t.Run("Success - With input", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			FindByNameFunc: func(ctx context.Context, name, lang string) ([]models.CountryFormResponse, error) {
				return []models.CountryFormResponse{response}, nil
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/country/search?name=bra",
			"", nil, nil,
		)
		handler.GetCountriesByName(c)

		expected := testutils.ToJSON(
			map[string][]models.CountryFormResponse{
				"data": {response},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Without input", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			FindAllFunc: func(ctx context.Context, lang string) ([]models.CountryFormResponse, error) {
				return []models.CountryFormResponse{response}, nil
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/country/search?name=",
			"", nil, nil,
		)
		handler.GetCountriesByName(c)

		expected := testutils.ToJSON(
			map[string][]models.CountryFormResponse{
				"data": {response},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockCountryService{
			FindByNameFunc: func(ctx context.Context, name, lang string) ([]models.CountryFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := country.NewAdminCountryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/country/search?name=bra",
			"", nil, nil,
		)
		handler.GetCountriesByName(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
