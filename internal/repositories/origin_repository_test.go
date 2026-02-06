package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
)

func TestNewOriginRepo(t *testing.T) {
	db := testutils.NewMockDB()
	originRepo := repositories.NewOriginRepo(db)

	assert.NotEmpty(t, originRepo)
}

func TestGetOrigins(t *testing.T) {
	db := testutils.NewMockDB()
	originRepo := repositories.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	db.Create(&origin)

	t.Run("Success", func(t *testing.T) {
		origins, err := originRepo.GetOrigins(context.Background())

		expected := []models.Origin{origin}

		assert.NoError(t, err)
		assert.Equal(t, expected, origins)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockOriginRepo := repositories.NewOriginRepo(mockDB)
		origins, err := mockOriginRepo.GetOrigins(context.Background())

		assert.Empty(t, origins)
		assert.Error(t, err)
	})
}

func TestGetActiveOrigins(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	origin2 := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"}, false)

	db.Create(&origin)
	db.Create(&origin2)

	t.Run("Success", func(t *testing.T) {
		origins, err := repo.GetActiveOrigins(context.Background())

		expected := []models.Origin{origin}

		assert.NoError(t, err)
		assert.Equal(t, expected, origins)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repositories.NewOriginRepo(mockDB)
		origins, err := mockCountryRepo.GetActiveOrigins(context.Background())

		assert.Empty(t, origins)
		assert.Error(t, err)
	})
}

func TestGetOriginByID(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.NewString(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)
	db.Create(&origin)

	t.Run("Success", func(t *testing.T) {
		resultOrigin, err := repo.GetOriginByID(context.Background(), origin.ID)

		assert.NoError(t, err)
		assert.Equal(t, origin, *resultOrigin)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repositories.NewOriginRepo(mockDB)
		origin, err := mockCountryRepo.GetOriginByID(context.Background(), uuid.UUID{})

		assert.Empty(t, origin)
		assert.Error(t, err)
	})
}

func TestGetOriginsByName(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewOriginRepo(db)

	origin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)
	db.Create(&origin)

	expected := []models.Origin{origin}

	t.Run("Success - Pt", func(t *testing.T) {
		lang := "pt"
		resultOrigins, err := repo.GetOriginsByName(context.Background(), "Alimentar", lang)

		assert.NoError(t, err)
		assert.Equal(t, expected, resultOrigins)
	})

	t.Run("Success - En", func(t *testing.T) {
		lang := "en"
		resultOrigins, err := repo.GetOriginsByName(context.Background(), "Food", lang)

		assert.NoError(t, err)
		assert.Equal(t, expected, resultOrigins)
	})

	t.Run("Success - Es", func(t *testing.T) {
		lang := "es"
		resultOrigins, err := repo.GetOriginsByName(context.Background(), "Alimentaria", lang)

		assert.NoError(t, err)
		assert.Equal(t, expected, resultOrigins)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repositories.NewOriginRepo(mockDB)
		origin, err := mockCountryRepo.GetOriginsByName(context.Background(), "", "")

		assert.Empty(t, origin)
		assert.Error(t, err)
	})
}

func TestGetOriginDuplicate(t *testing.T) {
	db := testutils.NewMockDB()
	originRepo := repositories.NewOriginRepo(db)

	mockOrigin := testmodels.NewOrigin(
		uuid.New().String(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)
	db.Create(&mockOrigin)

	t.Run("Success - With ID", func(t *testing.T) {
		origin, err := originRepo.GetOriginDuplicate(context.Background(), mockOrigin.Names, uuid.New())

		assert.NoError(t, err)
		assert.Equal(t, mockOrigin, *origin)
	})

	t.Run("Success - Without ID", func(t *testing.T) {
		origin, err := originRepo.GetOriginDuplicate(context.Background(), mockOrigin.Names, uuid.UUID{})

		assert.NoError(t, err)
		assert.Equal(t, mockOrigin, *origin)
	})

	t.Run("Error - Record not found", func(t *testing.T) {
		names := map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}
		origin, err := originRepo.GetOriginDuplicate(context.Background(), names, uuid.UUID{})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, origin)
	})

	t.Run("DB error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockOriginRepo := repositories.NewOriginRepo(mockDB)
		origin, err := mockOriginRepo.GetOriginDuplicate(context.Background(), mockOrigin.Names, uuid.New())

		assert.Empty(t, origin)
		assert.Error(t, err)
	})
}

func TestCreateOrigin(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)

	t.Run("Success", func(t *testing.T) {
		err := repo.CreateOrigin(context.Background(), &origin)
		assert.NoError(t, err)

		var result models.Origin
		err = db.Where("id = ?", origin.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, origin, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repositories.NewOriginRepo(mockDB)
		err = mockCountryRepo.CreateOrigin(context.Background(), &models.Origin{})

		assert.Error(t, err)
	})
}

func TestUpdateOrigin(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Hum", "en": "Human", "es": "Humanio"}, true)

	t.Run("Success", func(t *testing.T) {
		originToUpdate := models.Origin{
			ID:       origin.ID,
			Names:    map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
			IsActive: true,
		}

		err := repo.UpdateOrigin(context.Background(), &originToUpdate)
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

		mockCountryRepo := repositories.NewOriginRepo(mockDB)
		err = mockCountryRepo.UpdateOrigin(context.Background(), &models.Origin{})

		assert.Error(t, err)
	})
}

func TestDeleteOrigin(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewOriginRepo(db)

	origin := testmodels.NewOrigin(uuid.New().String(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)

	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteOrigin(context.Background(), &origin)
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

		mockCountryRepo := repositories.NewOriginRepo(mockDB)
		err = mockCountryRepo.DeleteOrigin(context.Background(), &models.Origin{})

		assert.Error(t, err)
	})
}
