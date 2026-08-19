package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewSampleRepo(t *testing.T) {
	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	assert.NotEmpty(t, sampleRepo)
}

func TestGetSamples(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := testmodels.CreateMockSample()
	db.Create(&mockSample)

	t.Run("Success - All samples", func(t *testing.T) {
		result, err := sampleRepo.GetSamples(ctx, "", uuid.Nil)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
	})

	t.Run("Success - Filtered samples", func(t *testing.T) {
		result, err := sampleRepo.GetSamples(ctx, mockSample.OriginCode, uuid.Nil)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
		assert.Equal(t, mockSample.OriginCode, result[0].OriginCode)
	})

	t.Run("Success - Filtered samples by user", func(t *testing.T) {
		result, err := sampleRepo.GetSamples(ctx, "", mockSample.UserID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		samples, err := mockSampleRepo.GetSamples(
			context.Background(), "", uuid.Nil)

		assert.Empty(t, samples)
		assert.Error(t, err)
	})
}

func TestGetSampleByID(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := testmodels.CreateMockSample()
	db.Create(&mockSample)

	t.Run("Success", func(t *testing.T) {
		resultSample, err := sampleRepo.GetSampleByID(ctx, mockSample.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockSample.ID, resultSample.ID)
		assert.Equal(t, mockSample.OriginCode, resultSample.OriginCode)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		resultSample, err := sampleRepo.GetSampleByID(ctx, uuid.New())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, resultSample)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		sample, err := mockSampleRepo.GetSampleByID(ctx, uuid.UUID{})

		assert.Empty(t, sample)
		assert.Error(t, err)
	})
}

func TestCreateSample(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := testmodels.CreateMockSample()

	t.Run("Success", func(t *testing.T) {
		err := sampleRepo.CreateSample(ctx, &mockSample)
		assert.NoError(t, err)

		var result models.Sample
		err = db.Where("id = ?", mockSample.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, mockSample.ID, result.ID)
		assert.Equal(t, mockSample.OriginCode, result.OriginCode)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		err = mockSampleRepo.CreateSample(ctx, &models.Sample{})

		assert.Error(t, err)
	})
}

func TestUpdateSample(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := testmodels.CreateMockSample()
	db.Create(&mockSample)

	t.Run("Success", func(t *testing.T) {
		sampleToUpdate := mockSample
		sampleToUpdate.OriginCode = "Updated Origin Code"

		err := sampleRepo.UpdateSample(ctx, &sampleToUpdate)
		assert.NoError(t, err)

		var result models.Sample
		err = db.Where("id = ?", mockSample.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, "Updated Origin Code", result.OriginCode)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		err = mockSampleRepo.UpdateSample(ctx, &models.Sample{})

		assert.Error(t, err)
	})
}

func TestDeleteSample(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	mockSample := testmodels.CreateMockSample()
	db.Create(&mockSample)

	t.Run("Success", func(t *testing.T) {
		err := sampleRepo.DeleteSample(ctx, &mockSample)
		assert.NoError(t, err)

		var result models.Sample
		err = db.Where("id = ?", mockSample.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		err = mockSampleRepo.DeleteSample(ctx, &models.Sample{})

		assert.Error(t, err)
	})
}
