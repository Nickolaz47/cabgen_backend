package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	country1 := testmodels.NewCountry("BRA", map[string]string{
		"pt": "Brasil", "en": "Brazil", "es": "Brazil",
	})
	country2 := testmodels.NewCountry("SPN", map[string]string{
		"pt": "Espanha", "en": "Spain", "es": "España",
	})
	countries := []models.Country{country1, country2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewCountrySeedRepository(db)

		err := repo.BulkInsert(context.Background(), countries)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewCountrySeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), countries)

		assert.Error(t, err)
	})
}

func TestCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewCountrySeedRepository(db)

	country1 := testmodels.NewCountry("BRA", map[string]string{
		"pt": "Brasil", "en": "Brazil", "es": "Brazil",
	})
	country2 := testmodels.NewCountry("SPN", map[string]string{
		"pt": "Espanha", "en": "Spain", "es": "España",
	})
	db.Create(&country1)
	db.Create(&country2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewCountrySeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}
