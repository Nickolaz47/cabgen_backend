package repository_test

import (
	"context"
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

func TestNewLaboratoryRepo(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	assert.NotEmpty(t, labRepo)
}

func TestGetLaboratories(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratory 1", "L1", true)
	db.Create(&lab)

	t.Run("Success", func(t *testing.T) {
		labs, err := labRepo.GetLaboratories(context.Background())

		expected := []models.Laboratory{lab}

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockLabRepo := repository.NewLaboratoryRepo(mockDB)
		labs, err := mockLabRepo.GetLaboratories(context.Background())

		assert.Empty(t, labs)
		assert.Error(t, err)
	})
}

func TestGetActiveLaboratories(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratory 1", "L1", true)
	lab2 := testmodels.NewLaboratory(uuid.NewString(), "Laboratory 2", "L2", false)
	db.Create(&lab)
	db.Create(&lab2)

	t.Run("Success", func(t *testing.T) {
		labs, err := labRepo.GetActiveLaboratories(context.Background())

		expected := []models.Laboratory{lab}

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockLabRepo := repository.NewLaboratoryRepo(mockDB)
		labs, err := mockLabRepo.GetActiveLaboratories(context.Background())

		assert.Empty(t, labs)
		assert.Error(t, err)
	})
}

func TestGetLaboratoryByID(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratory 1", "L1", true)
	db.Create(&lab)

	t.Run("Success", func(t *testing.T) {
		resultLab, err := labRepo.GetLaboratoryByID(context.Background(), lab.ID)

		assert.NoError(t, err)
		assert.Equal(t, lab, *resultLab)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockLabRepo := repository.NewLaboratoryRepo(mockDB)
		lab, err := mockLabRepo.GetLaboratoryByID(context.Background(), uuid.New())

		assert.Empty(t, lab)
		assert.Error(t, err)
	})
}

func TestGetLaboratoriesByNameOrAbbreviation(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Rio de Janeiro", "LACEN/RJ", true)
	lab2 := testmodels.NewLaboratory(uuid.NewString(), "Laboratorio Central do Para", "LACEN/PA", true)
	db.Create(&lab)
	db.Create(&lab2)

	t.Run("Success - Name", func(t *testing.T) {
		labs, err := labRepo.GetLaboratoriesByNameOrAbbreviation(context.Background(), "janeiro")

		expected := []models.Laboratory{lab}

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Success - Abbreviation", func(t *testing.T) {
		labs, err := labRepo.GetLaboratoriesByNameOrAbbreviation(context.Background(), "lacen")

		expected := []models.Laboratory{lab, lab2}

		assert.NoError(t, err)
		assert.Equal(t, expected, labs)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockLabRepo := repository.NewLaboratoryRepo(mockDB)
		labs, err := mockLabRepo.GetLaboratoriesByNameOrAbbreviation(context.Background(), "Lab")

		assert.Empty(t, labs)
		assert.Error(t, err)
	})
}

func TestCreateLaboratory(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratory 1", "L1", true)

	t.Run("Success", func(t *testing.T) {
		err := labRepo.CreateLaboratory(context.Background(), &lab)
		assert.NoError(t, err)

		var result models.Laboratory
		err = db.Where("id = ?", lab.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, lab, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockLabRepo := repository.NewLaboratoryRepo(mockDB)
		err = mockLabRepo.CreateLaboratory(context.Background(), &lab)

		assert.Error(t, err)
	})
}

func TestUpdateLaboratory(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratory 1", "L1", true)

	t.Run("Success", func(t *testing.T) {
		labToUpdate := models.Laboratory{
			ID:           lab.ID,
			Name:         lab.Name,
			Abbreviation: "Lab1",
			IsActive:     lab.IsActive,
		}

		err := labRepo.UpdateLaboratory(context.Background(), &labToUpdate)
		assert.NoError(t, err)

		var result models.Laboratory
		err = db.Where("id = ?", lab.ID).First(&result).Error

		expected := models.Laboratory{
			ID:           lab.ID,
			Name:         lab.Name,
			Abbreviation: labToUpdate.Abbreviation,
			IsActive:     lab.IsActive,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockLabRepo := repository.NewLaboratoryRepo(mockDB)
		err = mockLabRepo.UpdateLaboratory(context.Background(), &lab)

		assert.Error(t, err)
	})
}

func TestDeleteLaboratory(t *testing.T) {
	db := testutils.NewMockDB()
	labRepo := repository.NewLaboratoryRepo(db)

	lab := testmodels.NewLaboratory(uuid.NewString(), "Laboratory 1", "L1", true)

	t.Run("Success", func(t *testing.T) {
		err := labRepo.DeleteLaboratory(context.Background(), &lab)

		assert.NoError(t, err)

		var result models.Laboratory
		err = db.Where("id = ?", lab.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockLabRepo := repository.NewLaboratoryRepo(mockDB)
		err = mockLabRepo.DeleteLaboratory(context.Background(), &lab)

		assert.Error(t, err)
	})
}
