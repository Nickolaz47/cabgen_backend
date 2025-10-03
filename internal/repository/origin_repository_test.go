package repository_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
)

func TestNewOriginRepo(t *testing.T) {
	db := testutils.NewMockDB()
	result := repository.NewOriginRepo(db)

	assert.NotEmpty(t, result)
}

func TestGetOrigins(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	db.Create(&origin)
	t.Run("Success", func(t *testing.T) {
		origins, err := repo.GetOrigins()

		expected := []models.Origin{origin}

		assert.NoError(t, err)
		assert.NotEmpty(t, origins)
		assert.Equal(t, expected, origins)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewOriginRepo(mockDB)
		origins, err := mockCountryRepo.GetOrigins()

		assert.Empty(t, origins)
		assert.Error(t, err)
	})
}

func TestGetOriginByID(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewOriginRepo(db)

	id := uuid.New()
	origin := testmodels.NewOrigin(id.String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	db.Create(&origin)

	t.Run("Success", func(t *testing.T) {
		resultOrigin, err := repo.GetOriginByID(id)

		assert.NoError(t, err)
		assert.Equal(t, &origin, resultOrigin)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewOriginRepo(mockDB)
		origin, err := mockCountryRepo.GetOriginByID(uuid.UUID{})

		assert.Empty(t, origin)
		assert.Error(t, err)
	})
}

func TestGetOriginByName(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewOriginRepo(db)

	origin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)
	db.Create(&origin)
	t.Run("Success - Pt", func(t *testing.T) {
		lang := "pt"
		resultOrigin, err := repo.GetOriginByName("Alimentar", lang)

		assert.NoError(t, err)
		assert.Equal(t, &origin, resultOrigin)
	})

	t.Run("Success - En", func(t *testing.T) {
		lang := "en"
		resultOrigin, err := repo.GetOriginByName("Food", lang)

		assert.NoError(t, err)
		assert.Equal(t, &origin, resultOrigin)
	})

	t.Run("Success - Es", func(t *testing.T) {
		lang := "es"
		resultOrigin, err := repo.GetOriginByName("Alimentaria", lang)

		assert.NoError(t, err)
		assert.Equal(t, &origin, resultOrigin)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewOriginRepo(mockDB)
		origin, err := mockCountryRepo.GetOriginByName("", "")

		assert.Empty(t, origin)
		assert.Error(t, err)
	})
}

func TestCreateOrigin(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	t.Run("Success", func(t *testing.T) {
		err := repo.CreateOrigin(&origin)
		assert.NoError(t, err)

		var result models.Origin
		err = db.Where("id = ?", origin.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, origin, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewOriginRepo(mockDB)
		err = mockCountryRepo.CreateOrigin(&models.Origin{})

		assert.Error(t, err)
	})
}

func TestUpdateOrigin(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Hum", "en": "Human", "es": "Humanio"}, true)
	db.Create(&origin)
	t.Run("Success", func(t *testing.T) {
		originToUpdate := models.Origin{
			ID:       origin.ID,
			Names:    map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
			IsActive: true,
		}

		err := repo.UpdateOrigin(&originToUpdate)
		assert.NoError(t, err)

		var result models.Origin
		err = db.Where("id = ?", origin.ID).First(&result).Error

		expected := models.Origin{
			ID: origin.ID,
			Names: map[string]string{
				"pt": originToUpdate.Names["pt"], "en": "Human", "es": originToUpdate.Names["es"]},
			IsActive: true,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewOriginRepo(mockDB)
		err = mockCountryRepo.UpdateOrigin(&models.Origin{})

		assert.Error(t, err)
	})
}

func TestDeleteOrigin(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	db.Create(&origin)
	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteOrigin(&origin)

		assert.NoError(t, err)

		var result models.Origin
		err = db.Where("id = ?", origin.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewOriginRepo(mockDB)
		err = mockCountryRepo.DeleteOrigin(&models.Origin{})

		assert.Error(t, err)
	})
}
