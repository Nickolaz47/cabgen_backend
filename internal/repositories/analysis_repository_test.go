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

func TestNewAnalysisRepo(t *testing.T) {
	db := testutils.NewMockDB()
	result := repositories.NewAnalysisRepository(db)

	assert.NotEmpty(t, result)
}

func TestGetAnalyses(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	repo := repositories.NewAnalysisRepository(db)

	analysis := testmodels.CreateMockAnalysis()
	db.Create(&analysis)

	t.Run("Success", func(t *testing.T) {
		analyses, err := repo.GetAnalyses(ctx)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
		assert.Equal(t, analysis.ID, analyses[0].ID)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAnalysisRepo := repositories.NewAnalysisRepository(mockDB)
		analyses, err := mockAnalysisRepo.GetAnalyses(ctx)

		assert.Error(t, err)
		assert.Empty(t, analyses)
	})
}

func TestGetAnalysisByID(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	repo := repositories.NewAnalysisRepository(db)

	analysis := testmodels.CreateMockAnalysis()
	db.Create(&analysis)

	t.Run("Success", func(t *testing.T) {
		resultAnalysis, err := repo.GetAnalysisByID(ctx, analysis.ID)

		assert.NoError(t, err)
		assert.Equal(t, analysis.ID, resultAnalysis.ID)
		assert.Equal(t, analysis.Metrics, resultAnalysis.Metrics)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		resultAnalysis, err := repo.GetAnalysisByID(ctx, uuid.New())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, resultAnalysis)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAnalysisRepo := repositories.NewAnalysisRepository(mockDB)
		analysis, err := mockAnalysisRepo.GetAnalysisByID(ctx, uuid.UUID{})

		assert.Error(t, err)
		assert.Empty(t, analysis)
	})
}

func TestCreateAnalysis(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	repo := repositories.NewAnalysisRepository(db)

	analysis := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		err := repo.CreateAnalysis(ctx, &analysis)
		assert.NoError(t, err)

		var result models.Analysis
		err = db.Where("id = ?", analysis.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, analysis.ID, result.ID)
		assert.Equal(t, analysis.Metrics, result.Metrics)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAnalysisRepo := repositories.NewAnalysisRepository(mockDB)
		err = mockAnalysisRepo.CreateAnalysis(ctx, &models.Analysis{})

		assert.Error(t, err)
	})
}

func TestUpdateAnalysis(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	repo := repositories.NewAnalysisRepository(db)

	analysis := testmodels.CreateMockAnalysis()
	db.Create(&analysis)

	t.Run("Success", func(t *testing.T) {
		analysisToUpdate := analysis
		analysisToUpdate.Metrics = nil
		analysisToUpdate.StartedAt = nil
		analysisToUpdate.Status = models.AnalysisStatusPending

		err := repo.UpdateAnalysis(ctx, &analysisToUpdate)
		assert.NoError(t, err)

		var result models.Analysis
		err = db.Where("id = ?", analysis.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Nil(t, result.Metrics)
		assert.Nil(t, result.StartedAt)
		assert.Equal(t, result.Status, result.Status)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAnalysisRepo := repositories.NewAnalysisRepository(mockDB)
		err = mockAnalysisRepo.UpdateAnalysis(ctx, &models.Analysis{})

		assert.Error(t, err)
	})
}

func TestDeleteAnalysis(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	repo := repositories.NewAnalysisRepository(db)

	analysis := testmodels.CreateMockAnalysis()
	db.Create(&analysis)

	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteAnalysis(ctx, &analysis)
		assert.NoError(t, err)

		var result models.Analysis
		err = db.Where("id = ?", analysis.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAnalysisRepo := repositories.NewAnalysisRepository(mockDB)
		err = mockAnalysisRepo.DeleteAnalysis(ctx, &models.Analysis{})

		assert.Error(t, err)
	})
}
