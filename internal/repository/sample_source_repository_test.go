package repository_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewSampleSourceRepo(t *testing.T) {
	db := testutils.NewMockDB()
	result := repository.NewSampleSourceRepo(db)

	assert.NotEmpty(t, result)
}

func TestGetSampleSources(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSampleSourceRepo(db)

	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)
	db.Create(&sampleSource)

	t.Run("Success", func(t *testing.T) {
		sampleSources, err := repo.GetSampleSources()

		expected := []models.SampleSource{sampleSource}

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewSampleSourceRepo(mockDB)
		sampleSources, err := mockCountryRepo.GetSampleSources()

		assert.Empty(t, sampleSources)
		assert.Error(t, err)
	})
}

func TestGetActiveSampleSources(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSampleSourceRepo(db)

	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)
	sampleSource2 := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	db.Create(&sampleSource)
	db.Create(&sampleSource2)

	t.Run("Success", func(t *testing.T) {
		sampleSources, err := repo.GetActiveSampleSources()

		expected := []models.SampleSource{sampleSource}

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewSampleSourceRepo(mockDB)
		sampleSources, err := mockCountryRepo.GetActiveSampleSources()

		assert.Empty(t, sampleSources)
		assert.Error(t, err)
	})
}

func TestGetSampleSourceByID(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSampleSourceRepo(db)

	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado", "en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório", "en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)
	db.Create(&sampleSource)

	t.Run("Success", func(t *testing.T) {
		result, err := repo.GetSampleSourceByID(sampleSource.ID)

		assert.NoError(t, err)
		assert.Equal(t, &sampleSource, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewSampleSourceRepo(mockDB)
		result, err := mockCountryRepo.GetSampleSourceByID(uuid.UUID{})

		assert.Empty(t, result)
		assert.Error(t, err)
	})
}

func TestGetSampleSourcesByNameOrGroup(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSampleSourceRepo(db)

	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	sampleSource2 := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Coágulo sanguíneo", "en": "Blood clot", "es": "Coágulo de sangre"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	db.Create(&sampleSource)
	db.Create(&sampleSource2)

	t.Run("Success - Name", func(t *testing.T) {
		sampleSources, err := repo.GetSampleSourcesByNameOrGroup("plas", "en")

		expected := []models.SampleSource{sampleSource}

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Success - Group", func(t *testing.T) {
		sampleSources, err := repo.GetSampleSourcesByNameOrGroup("blo", "en")

		expected := []models.SampleSource{sampleSource, sampleSource2}

		assert.NoError(t, err)
		assert.Equal(t, expected, sampleSources)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewSampleSourceRepo(mockDB)
		sampleSources, err := mockCountryRepo.GetSampleSourcesByNameOrGroup("resp", "en")

		assert.Empty(t, sampleSources)
		assert.Error(t, err)
	})
}

func TestCreateSampleSource(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSampleSourceRepo(db)

	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)

	t.Run("Success", func(t *testing.T) {
		err := repo.CreateSampleSource(&sampleSource)
		assert.NoError(t, err)

		var result models.SampleSource
		err = db.Where("id = ?", sampleSource.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, sampleSource, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewSampleSourceRepo(mockDB)
		err = mockCountryRepo.CreateSampleSource(&sampleSource)

		assert.Error(t, err)
	})
}

func TestUpdateSampleSource(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSampleSourceRepo(db)

	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasm", "es": "Plasme"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	db.Create(&sampleSource)

	t.Run("Success", func(t *testing.T) {
		sampleSourceToUpdate := models.SampleSource{
			ID: sampleSource.ID,
			Names: map[string]string{
				"pt": "Plasma", "en": "Plasma", "es": "Plasma",
			},
			Groups:   sampleSource.Groups,
			IsActive: sampleSource.IsActive,
		}

		err := repo.UpdateSampleSource(&sampleSourceToUpdate)
		assert.NoError(t, err)

		var result models.SampleSource
		err = db.Where("id = ?", sampleSource.ID).First(&result).Error

		expected := models.SampleSource{
			ID:       sampleSource.ID,
			Names:    sampleSourceToUpdate.Names,
			Groups:   sampleSource.Groups,
			IsActive: sampleSource.IsActive,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewSampleSourceRepo(mockDB)
		err = mockCountryRepo.UpdateSampleSource(&models.SampleSource{})

		assert.Error(t, err)
	})
}

func TestDeleteSampleSource(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSampleSourceRepo(db)

	sampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasm", "es": "Plasme"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	db.Create(&sampleSource)

	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteSampleSource(&sampleSource)
		assert.NoError(t, err)

		var result models.SampleSource
		err = db.Where("id = ?", sampleSource.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockCountryRepo := repository.NewSampleSourceRepo(mockDB)
		err = mockCountryRepo.DeleteSampleSource(&models.SampleSource{})

		assert.Error(t, err)
	})
}
