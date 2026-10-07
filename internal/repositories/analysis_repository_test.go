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
	"gorm.io/datatypes"
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

func TestGetDashboardAnalyses(t *testing.T) {
	ctx := context.Background()

	newDashboardAnalysis := func(sampleID uuid.UUID,
		status models.AnalysisStatus, analysisType models.AnalysisType,
		createdAt time.Time) models.Analysis {
		analysis := testmodels.CreateMockAnalysis()
		analysis.Sample = models.Sample{}
		analysis.SampleID = sampleID
		analysis.Status = status
		analysis.Type = analysisType
		analysis.CreatedAt = createdAt
		return analysis
	}

	t.Run("Success - Latest DONE analysis per sample", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = true
		db.Create(&sample)

		older := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeGenome,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		db.Create(&older)

		newer := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeComplete,
			time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC))
		db.Create(&newer)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
		assert.Equal(t, newer.ID, analyses[0].ID)
		assert.Equal(t, sample.ID, analyses[0].Sample.ID)
		assert.Equal(t, "Aspirated",
			analyses[0].Sample.SampleSource.Names["en"])
	})

	t.Run("Success - Returns all in-network samples", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		firstSample := testmodels.CreateMockSample()
		firstSample.InNetwork = true
		db.Create(&firstSample)

		secondSample := testmodels.CreateMockSample()
		secondSample.InNetwork = true
		secondSample.HealthService.Name = "Outro Laboratorio Central"
		db.Create(&secondSample)

		first := newDashboardAnalysis(firstSample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeGenome,
			time.Date(2024, time.January, 3, 0, 0, 0, 0, time.UTC))
		db.Create(&first)

		second := newDashboardAnalysis(secondSample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeComplete,
			time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC))
		db.Create(&second)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Len(t, analyses, 2)
		assert.Equal(t, first.ID, analyses[0].ID)
		assert.Equal(t, second.ID, analyses[1].ID)
	})

	t.Run("Success - Empty database", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Empty(t, analyses)
	})

	t.Run("Success - Excludes sample outside network", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = false
		db.Create(&sample)

		analysis := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeGenome,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		db.Create(&analysis)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Empty(t, analyses)
	})

	t.Run("Success - Excludes statuses other than DONE", func(t *testing.T) {
		statuses := []models.AnalysisStatus{models.AnalysisStatusPending,
			models.AnalysisStatusRunning, models.AnalysisStatusFailed}

		for _, status := range statuses {
			t.Run(string(status), func(t *testing.T) {
				db := testutils.NewMockDB()
				repo := repositories.NewAnalysisRepository(db)

				sample := testmodels.CreateMockSample()
				sample.InNetwork = true
				db.Create(&sample)

				analysis := newDashboardAnalysis(sample.ID, status,
					models.AnalysisTypeGenome,
					time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
				db.Create(&analysis)

				analyses, err := repo.GetDashboardAnalyses(ctx)

				assert.NoError(t, err)
				assert.Empty(t, analyses)
			})
		}
	})

	t.Run("Success - Excludes FASTQC analyses", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = true
		db.Create(&sample)

		analysis := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeFastQC,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		db.Create(&analysis)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Empty(t, analyses)
	})

	t.Run("Success - Prefers DONE over newer non-DONE", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = true
		db.Create(&sample)

		done := newDashboardAnalysis(sample.ID, models.AnalysisStatusDone,
			models.AnalysisTypeGenome,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		db.Create(&done)

		failed := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusFailed, models.AnalysisTypeComplete,
			time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC))
		db.Create(&failed)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Len(t, analyses, 1)
		assert.Equal(t, done.ID, analyses[0].ID)
	})

	t.Run("Success - Excludes completeness below 95", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = true
		db.Create(&sample)

		analysis := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeGenome,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		analysis.Metrics = datatypes.JSON(
			`{"completeness": "80", "contamination": "1.0", "contigs": "100"}`)
		db.Create(&analysis)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Empty(t, analyses)
	})

	t.Run("Success - Excludes contamination above 5", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = true
		db.Create(&sample)

		analysis := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeGenome,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		analysis.Metrics = datatypes.JSON(
			`{"completeness": "99", "contamination": "9.9", "contigs": "100"}`)
		db.Create(&analysis)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Empty(t, analyses)
	})

	t.Run("Success - Excludes contigs above 500", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = true
		db.Create(&sample)

		analysis := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeGenome,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		analysis.Metrics = datatypes.JSON(
			`{"completeness": "99", "contamination": "1.0", "contigs": "600"}`)
		db.Create(&analysis)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Empty(t, analyses)
	})

	t.Run("Success - Excludes analyses without metrics", func(t *testing.T) {
		db := testutils.NewMockDB()
		repo := repositories.NewAnalysisRepository(db)

		sample := testmodels.CreateMockSample()
		sample.InNetwork = true
		db.Create(&sample)

		analysis := newDashboardAnalysis(sample.ID,
			models.AnalysisStatusDone, models.AnalysisTypeGenome,
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
		analysis.Metrics = nil
		db.Create(&analysis)

		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.NoError(t, err)
		assert.Empty(t, analyses)
	})

	t.Run("Error - No tables", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repo := repositories.NewAnalysisRepository(mockDB)
		analyses, err := repo.GetDashboardAnalyses(ctx)

		assert.Error(t, err)
		assert.Empty(t, analyses)
	})
}
