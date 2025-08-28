package country_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/country"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCountries(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", "", "", "")
	db.Create(&mockCountry)

	c, w := testutils.SetupGinContext(
		http.MethodGet, "/api/country", "",
		nil, nil,
	)

	country.GetCountries(c)
	expected := testutils.ToJSON(
		map[string][]models.Country{"data": {mockCountry}},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}

func TestGetCountryByID(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockCountry := testmodels.NewCountry("", "", "", "")
	db.Create(&mockCountry)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/country/BRA", "",
			nil, gin.Params{
				{Key: "code", Value: "BRA"},
			},
		)

		country.GetCountryByID(c)

		expected := testutils.ToJSON(
			map[string]models.Country{"data": mockCountry},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Country not found", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/country/XXX", "",
			nil, gin.Params{
				{Key: "code", Value: "XXX"},
			},
		)

		country.GetCountryByID(c)

		expected := `{"error": "No country was found with this code."}`

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
