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

func TestNewMicroorganismRepository(t *testing.T) {
	db := testutils.NewMockDB()
	result := repositories.NewMicroorganismRepository(db)

	assert.NotEmpty(t, result)
}

func TestGetMicroorganisms(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B",
			"en": "Serogroup B",
			"es": "Serogrupo B",
		},
		true,
	)
	db.Create(&micro)

	t.Run("Success", func(t *testing.T) {
		microorganisms, err := repo.GetMicroorganisms(
			context.Background())

		expected := []models.Microorganism{micro}

		assert.NoError(t, err)
		assert.Equal(t, expected, microorganisms)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		microorganisms, err := mockRepo.GetMicroorganisms(context.Background())

		assert.Empty(t, microorganisms)
		assert.Error(t, err)
	})
}

func TestGetActiveMicroorganisms(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	micro1 := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B",
			"en": "Serogroup B",
			"es": "Serogrupo B"},
		true,
	)
	micro2 := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Escherichia coli",
		map[string]string{
			"pt": "Genérico",
			"en": "Generic",
			"es": "Genérico"},
		false,
	)
	db.Create(&micro1)
	db.Create(&micro2)

	t.Run("Success", func(t *testing.T) {
		microorganisms, err := repo.GetActiveMicroorganisms(
			context.Background())

		expected := []models.Microorganism{micro1}

		assert.NoError(t, err)
		assert.Equal(t, expected, microorganisms)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		microorganisms, err := mockRepo.GetActiveMicroorganisms(
			context.Background())

		assert.Empty(t, microorganisms)
		assert.Error(t, err)
	})
}

func TestGetMicroorganismByID(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B",
			"en": "Serogroup B",
			"es": "Serogrupo B"},
		true,
	)
	db.Create(&micro)

	t.Run("Success", func(t *testing.T) {
		result, err := repo.GetMicroorganismByID(
			context.Background(), micro.ID)

		assert.NoError(t, err)
		assert.Equal(t, &micro, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		result, err := mockRepo.GetMicroorganismByID(
			context.Background(), uuid.UUID{})

		assert.Empty(t, result)
		assert.Error(t, err)
	})
}

func TestGetMicroorganismsBySpecies(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	micro1 := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B",
			"en": "Serogroup B",
			"es": "Serogrupo B"},
		true,
	)
	micro2 := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria gonorrhoeae",
		map[string]string{"pt": "Padrão", "en": "Standard", "es": "Estándar"},
		true,
	)
	db.Create(&micro1)
	db.Create(&micro2)

	t.Run("Success - Species", func(t *testing.T) {
		microorganisms, err := repo.GetMicroorganismsBySpecies(
			context.Background(), "meningitidis", "en")

		expected := []models.Microorganism{micro1}

		assert.NoError(t, err)
		assert.Equal(t, expected, microorganisms)
	})

	t.Run("Success - Variety", func(t *testing.T) {
		microorganisms, err := repo.GetMicroorganismsBySpecies(
			context.Background(), "stand", "en")

		expected := []models.Microorganism{micro2}

		assert.NoError(t, err)
		assert.Equal(t, expected, microorganisms)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		microorganisms, err := mockRepo.GetMicroorganismsBySpecies(
			context.Background(), "neisseria", "en")

		assert.Empty(t, microorganisms)
		assert.Error(t, err)
	})
}

func TestGetMicroorganismDuplicate(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	mockMicro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B",
			"en": "Serogroup B",
			"es": "Serogrupo B"},
		true,
	)
	db.Create(&mockMicro)

	mockSimpleMicro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Escherichia coli",
		nil,
		true,
	)
	db.Create(&mockSimpleMicro)

	t.Run("Success - With ID and Variety", func(t *testing.T) {
		micro, err := repo.GetMicroorganismDuplicate(
			context.Background(), mockMicro.Species,
			mockMicro.Variety, uuid.New(),
		)

		assert.NoError(t, err)
		assert.Equal(t, &mockMicro, micro)
	})

	t.Run("Success - Without ID and Variety", func(t *testing.T) {
		micro, err := repo.GetMicroorganismDuplicate(
			context.Background(), mockMicro.Species,
			mockMicro.Variety, uuid.UUID{},
		)

		assert.NoError(t, err)
		assert.Equal(t, &mockMicro, micro)
	})

	t.Run("Success - Variety Nil", func(t *testing.T) {
		micro, err := repo.GetMicroorganismDuplicate(
			context.Background(), "Escherichia coli", nil, 
			uuid.New(),
		)

		assert.NoError(t, err)
		assert.Equal(t, mockSimpleMicro.ID, micro.ID)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		variety := map[string]string{
			"pt": "Outro", "en": "Other", "es": "Otro"}
		micro, err := repo.GetMicroorganismDuplicate(
			context.Background(), "Outra Specie", variety, uuid.UUID{},
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, micro)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		micro, err := mockRepo.GetMicroorganismDuplicate(
			context.Background(), mockMicro.Species,
			mockMicro.Variety, uuid.New(),
		)

		assert.Error(t, err)
		assert.Empty(t, micro)
	})
}

func TestCreateMicroorganism(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B",
			"en": "Serogroup B",
			"es": "Serogrupo B"},
		true,
	)

	t.Run("Success", func(t *testing.T) {
		err := repo.CreateMicroorganism(context.Background(), &micro)
		assert.NoError(t, err)

		var result models.Microorganism
		err = db.Where("id = ?", micro.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, micro, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		err = mockRepo.CreateMicroorganism(context.Background(), &micro)

		assert.Error(t, err)
	})
}

func TestUpdateMicroorganism(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo A",
			"en": "Serogroup A",
			"es": "Serogrupo A"},
		true,
	)
	db.Create(&micro)

	t.Run("Success", func(t *testing.T) {
		microToUpdate := models.Microorganism{
			ID:      micro.ID,
			Taxon:   micro.Taxon,
			Species: "Neisseria Updated",
			Variety: map[string]string{
				"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B",
			},
			IsActive: false,
		}

		err := repo.UpdateMicroorganism(context.Background(), &microToUpdate)
		assert.NoError(t, err)

		var result models.Microorganism
		err = db.Where("id = ?", micro.ID).First(&result).Error

		expected := models.Microorganism{
			ID:       micro.ID,
			Taxon:    microToUpdate.Taxon,
			Species:  microToUpdate.Species,
			Variety:  microToUpdate.Variety,
			IsActive: microToUpdate.IsActive,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		err = mockRepo.UpdateMicroorganism(context.Background(),
			&models.Microorganism{})

		assert.Error(t, err)
	})
}

func TestDeleteMicroorganism(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewMicroorganismRepository(db)

	micro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B",
			"en": "Serogroup B",
			"es": "Serogrupo B"},
		true,
	)
	db.Create(&micro)

	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteMicroorganism(context.Background(), &micro)
		assert.NoError(t, err)

		var result models.Microorganism
		err = db.Where("id = ?", micro.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewMicroorganismRepository(mockDB)
		err = mockRepo.DeleteMicroorganism(context.Background(),
			&models.Microorganism{})

		assert.Error(t, err)
	})
}
