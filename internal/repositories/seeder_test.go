package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCountryBulkInsert(t *testing.T) {
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

func TestCountryCount(t *testing.T) {
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

func TestMicroorganismBulkInsert(t *testing.T) {
	db := testutils.NewMockDB()

	micro1 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Escherichia coli",
		map[string]string{
			"pt": "Enteropatogênica A",
			"en": "Enteropathogenic A",
			"es": "Enteropatogénica A",
		}, true,
	)
	micro2 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Plesiomonas shigelloides",
		nil, true,
	)
	micros := []models.Microorganism{micro1, micro2}

	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewMicroorganismSeedRepository(db)

		err := repo.BulkInsert(context.Background(), micros)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewMicroorganismSeedRepository(mockDB)
		err = repo.BulkInsert(context.Background(), micros)

		assert.Error(t, err)
	})
}

func TestMicroorganismCount(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismSeedRepository(db)

	micro1 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Escherichia coli",
		map[string]string{
			"pt": "Enteropatogênica A",
			"en": "Enteropathogenic A",
			"es": "Enteropatogénica A",
		}, true,
	)
	micro2 := testmodels.NewMicroorganism(
		uuid.NewString(), models.Bacteria, "Plesiomonas shigelloides",
		nil, true,
	)
	db.Create(&micro1)
	db.Create(&micro2)

	t.Run("Success", func(t *testing.T) {
		count, err := repo.Count(context.Background())

		var expected int64 = 2

		assert.NoError(t, err)
		assert.Equal(t, expected, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismSeedRepository(mockDB)
		count, err := mockRepo.Count(context.Background())

		assert.Error(t, err)
		assert.Empty(t, count)
	})
}
