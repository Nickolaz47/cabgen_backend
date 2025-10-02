package repository_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
)

func TestNewCountryRepo(t *testing.T) {
	db := testutils.NewMockDB()

	countryRepo := repository.NewCountryRepo(db)

	assert.NotEmpty(t, countryRepo)
}

func TestGetCountries(t *testing.T) {
	db := testutils.NewMockDB()
	mockCountry := testmodels.NewCountry("", nil)
	mockCountry2 := testmodels.NewCountry(
		"CYP", map[string]string{"pt": "Chipre", "en": "Cyprus", "es": "Chipre"})
	db.Create(&mockCountry)
	db.Create(&mockCountry2)

	countryRepo := repository.NewCountryRepo(db)

	t.Run("Success", func(t *testing.T) {
		countries, err := countryRepo.GetCountries()

		expected := []models.Country{mockCountry, mockCountry2}

		assert.NoError(t, err)
		assert.Equal(t, expected, countries)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewCountryRepo(mockDB)
		countries, err := mockCountryRepo.GetCountries()

		assert.Empty(t, countries)
		assert.Error(t, err)
	})
}

func TestGetCountry(t *testing.T) {
	db := testutils.NewMockDB()
	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	countryRepo := repository.NewCountryRepo(db)

	t.Run("Success", func(t *testing.T) {
		country, err := countryRepo.GetCountry("BRA")
		expected := models.Country{
			Code: "BRA",
			Names: map[string]string{
				"en": "Brazil",
				"es": "Brazil",
				"pt": "Brasil",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, &expected, country)
	})

	t.Run("Error", func(t *testing.T) {
		country, err := countryRepo.GetCountry("ARG")

		assert.Empty(t, country)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
	})
}
