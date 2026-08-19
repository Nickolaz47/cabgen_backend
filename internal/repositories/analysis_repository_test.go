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

	t.Run("Success - userID is nil", func(t *testing.T) {
		analyses, err := repo.GetAnalyses(ctx, uuid.Nil, models.AnalysisFilter{})

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
		assert.Equal(t, analysis.ID, analyses[0].ID)
	})

	t.Run("Success - userID filter", func(t *testing.T) {
		analyses, err := repo.GetAnalyses(ctx, uuid.New(), models.AnalysisFilter{})

		assert.NoError(t, err)
		assert.Len(t, analyses, 0)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAnalysisRepo := repositories.NewAnalysisRepository(mockDB)
		analyses, err := mockAnalysisRepo.GetAnalyses(ctx, uuid.Nil, models.AnalysisFilter{})

		assert.Error(t, err)
		assert.Empty(t, analyses)
	})
}

func TestGetAnalysesFilters(t *testing.T) {
	ctx := context.Background()

	mockUser := testmodels.NewLoginUser()
	mockAnalysis := testmodels.CreateMockAnalysis()

	t.Run("Collaborator - Filter by Type", func(t *testing.T) {
		filterDB := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(filterDB)

		filterDB.Create(&mockAnalysis)

		filter := models.AnalysisFilter{Type: models.AnalysisTypeComplete}
		analyses, err := repo.GetAnalyses(ctx, mockUser.ID, filter)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
	})

	t.Run("Admin - Filter by Type", func(t *testing.T) {
		filterDB := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(filterDB)

		filterDB.Create(&mockAnalysis)

		filter := models.AnalysisFilter{Type: models.AnalysisTypeComplete}
		analyses, err := repo.GetAnalyses(ctx, uuid.Nil, filter)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
	})

	t.Run("Admin - Filter by Username", func(t *testing.T) {
		filterDB := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(filterDB)

		filterDB.Create(&mockAnalysis)

		filter := models.AnalysisFilter{Username: mockUser.Username}
		analyses, err := repo.GetAnalyses(ctx, uuid.Nil, filter)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
	})

	t.Run("Collaborator - Filter by Type Returns Subset", func(t *testing.T) {
		filterDB := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(filterDB)

		genomeAnalysis := mockAnalysis
		genomeAnalysis.ID = uuid.New()
		genomeAnalysis.Type = models.AnalysisTypeGenome
		filterDB.Create(&genomeAnalysis)

		completeAnalysis := mockAnalysis
		completeAnalysis.ID = uuid.New()
		completeAnalysis.Type = models.AnalysisTypeComplete
		filterDB.Create(&completeAnalysis)

		filter := models.AnalysisFilter{Type: models.AnalysisTypeGenome}
		analyses, err := repo.GetAnalyses(ctx, mockUser.ID, filter)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
		assert.Equal(t, genomeAnalysis.ID, analyses[0].ID)
	})

	t.Run("Collaborator - Filter by OriginCode", func(t *testing.T) {
		filterDB := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(filterDB)

		filterDB.Create(&mockAnalysis)

		filter := models.AnalysisFilter{
			OriginCode: mockAnalysis.Sample.OriginCode,
		}
		analyses, err := repo.GetAnalyses(ctx, mockUser.ID, filter)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
	})

	t.Run("Admin - Filter by OriginCode", func(t *testing.T) {
		filterDB := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(filterDB)

		filterDB.Create(&mockAnalysis)

		filter := models.AnalysisFilter{
			OriginCode: mockAnalysis.Sample.OriginCode,
		}
		analyses, err := repo.GetAnalyses(ctx, uuid.Nil, filter)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
	})

	t.Run("Empty filter - both paths", func(t *testing.T) {
		filterDB := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(filterDB)

		filterDB.Create(&mockAnalysis)

		analyses, err := repo.GetAnalyses(ctx, mockUser.ID, models.AnalysisFilter{})
		assert.NoError(t, err)
		assert.Len(t, analyses, 1)

		analyses, err = repo.GetAnalyses(ctx, uuid.Nil, models.AnalysisFilter{})
		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
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

func TestAnalysisUpdateSample(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	repo := repositories.NewAnalysisRepository(db)

	analysis := testmodels.CreateMockAnalysis()
	db.Create(&analysis)

	t.Run("Success", func(t *testing.T) {
		sampleToUpdate := analysis.Sample
		fasta := "/new/path/assembly.fasta"
		sampleToUpdate.Fasta = &fasta

		err := repo.UpdateSample(ctx, &sampleToUpdate)
		assert.NoError(t, err)

		var result models.Sample
		err = db.Where("id = ?", analysis.Sample.ID).First(&result).Error

		assert.NoError(t, err)
		assert.NotNil(t, result.Fasta)
		assert.Equal(t, fasta, *result.Fasta)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAnalysisRepo := repositories.NewAnalysisRepository(mockDB)
		err = mockAnalysisRepo.UpdateSample(ctx, &models.Sample{})

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
