package repositories_test

import (
	"context"
	"testing"
	"time"

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
		result, _, err := sampleRepo.GetSamples(ctx, "", uuid.Nil, 0, 0)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
	})

	t.Run("Success - Filtered samples", func(t *testing.T) {
		result, _, err := sampleRepo.GetSamples(ctx, mockSample.OriginCode, uuid.Nil, 0, 0)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
		assert.Equal(t, mockSample.OriginCode, result[0].OriginCode)
	})

	t.Run("Success - Filtered samples by user", func(t *testing.T) {
		result, _, err := sampleRepo.GetSamples(ctx, "", mockSample.UserID, 0, 0)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockSample.ID, result[0].ID)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockSampleRepo := repositories.NewSampleRepo(mockDB)
		samples, _, err := mockSampleRepo.GetSamples(
			context.Background(), "", uuid.Nil, 0, 0)

		assert.Empty(t, samples)
		assert.Error(t, err)
	})

	t.Run("Success - Pagination", func(t *testing.T) {
		pagDB := testutils.NewMockDB()
		pagRepo := repositories.NewSampleRepo(pagDB)

		first := testmodels.CreateMockSample()
		pagDB.Create(&first)

		second := testmodels.CreateMockSample()
		second.HealthService.Name = "Outro Laboratorio Central"
		pagDB.Create(&second)

		third := testmodels.CreateMockSample()
		third.HealthService.Name = "Terceiro Laboratorio Central"
		pagDB.Create(&third)

		samples, total, err := pagRepo.GetSamples(ctx, "", uuid.Nil, 2, 0)

		assert.NoError(t, err)
		assert.Len(t, samples, 2)
		assert.Equal(t, int64(3), total)

		samples, total, err = pagRepo.GetSamples(ctx, "", uuid.Nil, 2, 2)

		assert.NoError(t, err)
		assert.Len(t, samples, 1)
		assert.Equal(t, int64(3), total)

		samples, total, err = pagRepo.GetSamples(ctx, "", uuid.Nil, 0, 0)

		assert.NoError(t, err)
		assert.Len(t, samples, 3)
		assert.Equal(t, int64(0), total)

		samples, total, err = pagRepo.GetSamples(ctx, "", uuid.Nil, 2, 10)

		assert.NoError(t, err)
		assert.Empty(t, samples)
		assert.Equal(t, int64(3), total)
	})

	t.Run("Success - Ordered newest first", func(t *testing.T) {
		orderDB := testutils.NewMockDB()
		orderRepo := repositories.NewSampleRepo(orderDB)

		oldest := testmodels.CreateMockSample()
		oldest.CreatedAt = time.Date(2020, time.January, 1,
			0, 0, 0, 0, time.UTC)
		orderDB.Create(&oldest)

		middle := testmodels.CreateMockSample()
		middle.HealthService.Name = "Outro Laboratorio Central"
		middle.CreatedAt = time.Date(2022, time.January, 1,
			0, 0, 0, 0, time.UTC)
		orderDB.Create(&middle)

		newest := testmodels.CreateMockSample()
		newest.HealthService.Name = "Terceiro Laboratorio Central"
		newest.CreatedAt = time.Date(2024, time.January, 1,
			0, 0, 0, 0, time.UTC)
		orderDB.Create(&newest)

		samples, _, err := orderRepo.GetSamples(ctx, "", uuid.Nil, 2, 0)

		assert.NoError(t, err)
		assert.Len(t, samples, 2)
		assert.Equal(t, newest.ID, samples[0].ID)
		assert.Equal(t, middle.ID, samples[1].ID)

		samples, _, err = orderRepo.GetSamples(ctx, "", uuid.Nil, 0, 0)

		assert.NoError(t, err)
		assert.Len(t, samples, 3)
		assert.Equal(t, newest.ID, samples[0].ID)
		assert.Equal(t, oldest.ID, samples[2].ID)
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

func TestDeleteSampleCascadeAnalyses(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	sampleRepo := repositories.NewSampleRepo(db)

	analysis := testmodels.CreateMockAnalysis()
	db.Create(&analysis)

	t.Run("Success", func(t *testing.T) {
		err := sampleRepo.DeleteSample(ctx, &analysis.Sample)

		assert.NoError(t, err)

		var result models.Analysis
		err = db.Where("id = ?", analysis.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})
}

func TestCreateSamples(t *testing.T) {
	stripAssocs := func(s models.Sample) models.Sample {
		s.Country = models.Country{}
		s.User = models.User{}
		s.Origin = models.Origin{}
		s.SampleSource = models.SampleSource{}
		s.Microorganism = models.Microorganism{}
		s.Sequencer = models.Sequencer{}
		s.Laboratory = models.Laboratory{}
		s.HealthService = models.HealthService{}
		return s
	}

	countSamples := func(db *gorm.DB) int64 {
		var count int64
		db.Model(&models.Sample{}).Count(&count)
		return count
	}

	t.Run("Success", func(t *testing.T) {
		db := testutils.NewMockDB()
		sampleRepo := repositories.NewSampleRepo(db)

		base := testmodels.CreateMockSample()
		assert.NoError(t, db.Create(&base).Error)

		a := stripAssocs(base)
		a.ID = uuid.New()
		b := stripAssocs(base)
		b.ID = uuid.New()

		err := sampleRepo.CreateSamples(context.Background(),
			[]models.Sample{a, b})

		assert.NoError(t, err)
		assert.Equal(t, int64(3), countSamples(db))
	})

	t.Run("Rollback - duplicate id", func(t *testing.T) {
		db := testutils.NewMockDB()
		sampleRepo := repositories.NewSampleRepo(db)

		base := testmodels.CreateMockSample()
		assert.NoError(t, db.Create(&base).Error)

		a := stripAssocs(base)
		a.ID = uuid.New()
		b := stripAssocs(base)
		b.ID = a.ID

		err := sampleRepo.CreateSamples(context.Background(),
			[]models.Sample{a, b})

		assert.Error(t, err)
		assert.Equal(t, int64(1), countSamples(db))
	})
}
