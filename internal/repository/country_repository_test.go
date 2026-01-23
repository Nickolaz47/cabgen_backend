package repository_test

import (
	"context"
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
	countryRepo := repository.NewCountryRepo(db)

	mockCountry := testmodels.NewCountry("", nil)
	mockCountry2 := testmodels.NewCountry(
		"CYP", map[string]string{"pt": "Chipre", "en": "Cyprus", "es": "Chipre"})
	db.Create(&mockCountry)
	db.Create(&mockCountry2)

	t.Run("Success", func(t *testing.T) {
		countries, err := countryRepo.GetCountries(context.Background())

		expected := []models.Country{mockCountry, mockCountry2}

		assert.NoError(t, err)
		assert.Equal(t, expected, countries)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewCountryRepo(mockDB)
		countries, err := mockCountryRepo.GetCountries(context.Background())

		assert.Empty(t, countries)
		assert.Error(t, err)
	})
}

func TestGetCountryByID(t *testing.T) {
	db := testutils.NewMockDB()
	countryRepo := repository.NewCountryRepo(db)

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	t.Run("Success", func(t *testing.T) {
		country, err := countryRepo.GetCountryByID(context.Background(), mockCountry.ID)
		expected := models.Country{
			ID:   1,
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

	t.Run("Error - Not Found", func(t *testing.T) {
		country, err := countryRepo.GetCountryByID(context.Background(), 2)

		assert.Empty(t, country)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
	})

	t.Run("Error - DB", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewCountryRepo(mockDB)
		countries, err := mockCountryRepo.GetCountryByID(context.Background(), mockCountry.ID)

		assert.Empty(t, countries)
		assert.Error(t, err)
	})
}

func TestGetCountryByCode(t *testing.T) {
	db := testutils.NewMockDB()
	countryRepo := repository.NewCountryRepo(db)

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	t.Run("Success", func(t *testing.T) {
		country, err := countryRepo.GetCountryByCode(context.Background(), "BRA")
		expected := models.Country{
			ID:   1,
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

	t.Run("Error - Not Found", func(t *testing.T) {
		country, err := countryRepo.GetCountryByCode(context.Background(), "ARG")

		assert.Empty(t, country)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
	})

	t.Run("Error - DB", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewCountryRepo(mockDB)
		countries, err := mockCountryRepo.GetCountryByCode(context.Background(), "BRA")

		assert.Empty(t, countries)
		assert.Error(t, err)
	})
}

func TestGetCountriesByName(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewCountryRepo(db)

	country := testmodels.NewCountry("", nil)
	db.Create(&country)

	expected := []models.Country{country}

	t.Run("Success - Pt", func(t *testing.T) {
		lang := "pt"
		countries, err := repo.GetCountriesByName(context.Background(), "brasil", lang)

		assert.NoError(t, err)
		assert.Equal(t, expected, countries)
	})

	t.Run("Success - En", func(t *testing.T) {
		lang := "en"
		countries, err := repo.GetCountriesByName(context.Background(), "brazil", lang)

		assert.NoError(t, err)
		assert.Equal(t, expected, countries)
	})

	t.Run("Success - Es", func(t *testing.T) {
		lang := "es"
		countries, err := repo.GetCountriesByName(context.Background(), "brazil", lang)

		assert.NoError(t, err)
		assert.Equal(t, expected, countries)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewCountryRepo(mockDB)
		countries, err := mockCountryRepo.GetCountriesByName(context.Background(), "bra", "pt")

		assert.Error(t, err)
		assert.Empty(t, countries)
	})
}

func TestGetCountryDuplicate(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewCountryRepo(db)

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	t.Run("Success - With Code", func(t *testing.T) {
		country, err := repo.GetCountryDuplicate(
			context.Background(), mockCountry.Names,
			"ARG",
		)

		assert.NoError(t, err)
		assert.Equal(t, mockCountry, *country)
	})

	t.Run("Success - Without Code", func(t *testing.T) {
		country, err := repo.GetCountryDuplicate(
			context.Background(), mockCountry.Names,
			"",
		)

		assert.NoError(t, err)
		assert.Equal(t, mockCountry, *country)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		names := map[string]string{
			"pt": "Espanha",
			"en": "Spain",
			"es": "España",
		}
		country, err := repo.GetCountryDuplicate(
			context.Background(),
			names, "",
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, country)
	})

	t.Run("Error - DB", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewCountryRepo(mockDB)
		country, err := mockCountryRepo.GetCountryDuplicate(
			context.Background(), mockCountry.Names, "",
		)

		assert.Error(t, err)
		assert.Empty(t, country)
	})
}

func TestCreateCountry(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewCountryRepo(db)

	country := testmodels.NewCountry("", nil)

	t.Run("Success", func(t *testing.T) {
		err := repo.CreateCountry(context.Background(), &country)
		assert.NoError(t, err)

		var result models.Country
		err = db.Where("code = ?", country.Code).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, country, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"))
		assert.NoError(t, err)

		mockRepo := repository.NewCountryRepo(mockDB)
		err = mockRepo.CreateCountry(context.Background(), &country)

		assert.Error(t, err)
	})
}

func TestUpdateCountry(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewCountryRepo(db)

	country := testmodels.NewCountry(
		"SPN",
		map[string]string{
			"pt": "Espaha",
			"en": "Spain",
			"es": "España",
		},
	)

	t.Run("Success", func(t *testing.T) {
		countryToUpdate := models.Country{
			Code: country.Code,
			Names: map[string]string{
				"pt": "Espanha",
				"en": "Spain",
				"es": "España",
			},
		}

		err := repo.UpdateCountry(
			context.Background(),
			&countryToUpdate,
		)
		assert.NoError(t, err)

		var result models.Country
		err = db.Where("code = ?", country.Code).First(&result).Error

		expected := models.Country{
			ID:    1,
			Code:  country.Code,
			Names: countryToUpdate.Names,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repository.NewCountryRepo(mockDB)
		err = mockRepo.UpdateCountry(context.Background(), &models.Country{})

		assert.Error(t, err)
	})
}

func TestDeleteCountry(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewCountryRepo(db)

	country := testmodels.NewCountry("", nil)
	country.ID = 1

	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteCountry(context.Background(), &country)
		assert.NoError(t, err)

		var result models.Country
		err = db.Where("code = ?", country.Code).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repository.NewCountryRepo(mockDB)
		err = mockRepo.DeleteCountry(context.Background(), &models.Country{})

		assert.Error(t, err)
	})
}
