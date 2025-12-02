package repository_test

import (
	"context"
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

func TestNewSequencerRepo(t *testing.T) {
	db := testutils.NewMockDB()
	result := repository.NewSequencerRepo(db)

	assert.NotEmpty(t, result)
}

func TestGetSequencers(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)
	db.Create(&sequencer)
	t.Run("Success", func(t *testing.T) {
		sequencers, err := repo.GetSequencers(context.Background())

		expected := []models.Sequencer{sequencer}

		assert.NoError(t, err)
		assert.Equal(t, expected, sequencers)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		sequencers, err := mockSequencerRepo.GetSequencers(context.Background())

		assert.Error(t, err)
		assert.Empty(t, sequencers)
	})
}

func TestGetActiveSequencers(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)
	sequencer2 := testmodels.NewSequencer(uuid.NewString(), "Nanopore", "MinION", false)
	db.Create(&sequencer)
	db.Create(&sequencer2)

	t.Run("Success", func(t *testing.T) {
		sequencers, err := repo.GetActiveSequencers(context.Background())

		expected := []models.Sequencer{sequencer}

		assert.NoError(t, err)
		assert.Equal(t, expected, sequencers)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		sequencers, err := mockSequencerRepo.GetActiveSequencers(context.Background())

		assert.Error(t, err)
		assert.Empty(t, sequencers)
	})
}

func TestGetSequencerByID(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	id := uuid.New()
	sequencer := testmodels.NewSequencer(id.String(), "Illumina", "MiSeq", true)
	db.Create(&sequencer)

	t.Run("Success", func(t *testing.T) {
		resultSequencer, err := repo.GetSequencerByID(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, &sequencer, resultSequencer)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		sequencer, err := mockSequencerRepo.GetSequencerByID(context.Background(), uuid.UUID{})

		assert.Error(t, err)
		assert.Empty(t, sequencer)
	})
}

func TestGetSequencersByBrandOrModel(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)
	db.Create(&sequencer)
	sequencer2 := testmodels.NewSequencer(uuid.NewString(), "Illumina", "NextSeq 550", true)
	db.Create(&sequencer2)

	t.Run("Success - Brand", func(t *testing.T) {
		resultSequencer, err := repo.GetSequencersByBrandOrModel(context.Background(), "illu")

		expected := []models.Sequencer{sequencer, sequencer2}

		assert.NoError(t, err)
		assert.Equal(t, expected, resultSequencer)
	})

	t.Run("Success - Model", func(t *testing.T) {
		resultSequencer, err := repo.GetSequencersByBrandOrModel(context.Background(), "mis")

		expected := []models.Sequencer{sequencer}

		assert.NoError(t, err)
		assert.Equal(t, expected, resultSequencer)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		sequencers, err := mockSequencerRepo.GetSequencersByBrandOrModel(context.Background(), "illumina")

		assert.Error(t, err)
		assert.Empty(t, sequencers)
	})
}

func TestGetSequencerDuplicate(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	mockSequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)
	db.Create(&mockSequencer)

	t.Run("Success - With ID", func(t *testing.T) {
		sequencer, err := repo.GetSequencerDuplicate(context.Background(), mockSequencer.Model, uuid.New())

		assert.NoError(t, err)
		assert.Equal(t, &mockSequencer, sequencer)
	})

	t.Run("Success - Without ID", func(t *testing.T) {
		sequencer, err := repo.GetSequencerDuplicate(context.Background(), mockSequencer.Model, uuid.UUID{})

		assert.NoError(t, err)
		assert.Equal(t, &mockSequencer, sequencer)
	})

	t.Run("Error - Record not found", func(t *testing.T) {
		model := "MinION"
		sequencer, err := repo.GetSequencerDuplicate(context.Background(), model, uuid.UUID{})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, sequencer)
	})

	t.Run("DB error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		sequencer, err := mockSequencerRepo.GetSequencerDuplicate(context.Background(), "Mis", uuid.UUID{})

		assert.Error(t, err)
		assert.Empty(t, sequencer)
	})
}

func TestCreateSequencer(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MiSeq", true)
	t.Run("Success", func(t *testing.T) {
		err := repo.CreateSequencer(context.Background(), &sequencer)
		assert.NoError(t, err)

		var result models.Sequencer
		err = db.Where("id = ?", sequencer.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, sequencer, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		err = mockSequencerRepo.CreateSequencer(context.Background(), &models.Sequencer{})

		assert.Error(t, err)
	})
}

func TestUpdateSequencer(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MySeq", true)
	db.Create(&sequencer)
	t.Run("Success", func(t *testing.T) {
		sequencerToUpdate := models.Sequencer{
			ID:       sequencer.ID,
			Brand:    sequencer.Brand,
			Model:    "MiSeq",
			IsActive: true,
		}

		err := repo.UpdateSequencer(context.Background(), &sequencerToUpdate)
		assert.NoError(t, err)

		var result models.Sequencer
		err = db.Where("id = ?", sequencer.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, sequencerToUpdate, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		err = mockSequencerRepo.UpdateSequencer(context.Background(), &models.Sequencer{})

		assert.Error(t, err)
	})
}

func TestDeleteSequencer(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repository.NewSequencerRepo(db)

	sequencer := testmodels.NewSequencer(uuid.NewString(), "Illumina", "MySeq", true)
	db.Create(&sequencer)
	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteSequencer(context.Background(), &sequencer)
		assert.NoError(t, err)

		var result models.Sequencer
		err = db.Where("id = ?", sequencer.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSequencerRepo := repository.NewSequencerRepo(mockDB)
		err = mockSequencerRepo.DeleteSequencer(context.Background(), &models.Sequencer{})

		assert.Error(t, err)
	})
}
