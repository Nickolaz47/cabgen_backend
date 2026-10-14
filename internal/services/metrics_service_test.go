package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestMetricsGetMetrics(t *testing.T) {
	ctx := context.Background()

	mockSample := testmodels.CreateMockSample()
	mockUser := testmodels.NewLoginUser()

	var mockCountryCode = "BRA"

	t.Run("Success", func(t *testing.T) {
		mockDone := testmodels.CreateMockAnalysis()
		mockDone.Status = models.AnalysisStatusDone
		mockDone.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","acquired_resistance":["blaOXA-23"]}`)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string, userID uuid.UUID, limit, offset int) ([]models.Sample, int64, error) {
				return []models.Sample{mockSample}, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID, filter models.AnalysisFilter, limit, offset int) ([]models.Analysis, int64, error) {
				return []models.Analysis{mockDone}, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context, filter models.AdminUserFilter, limit, offset int) ([]models.User, int64, error) {
				return []models.User{mockUser}, 0, nil
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo, nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.TotalSamples)
		assert.Equal(t, int64(1), result.TotalCountries)
		assert.Equal(t, int64(1), result.TotalUsers)
		assert.Equal(t, int64(1), result.TotalAnalyses)
		assert.Equal(t, int64(1), result.AnalysesByStatus.Done)
		assert.Equal(t, int64(1), result.TotalSpecies)
		assert.Equal(t, int64(1), result.TotalResistance)
		assert.Equal(t, []models.CountryMetric{{Country: mockCountryCode, Count: 1}},
			result.TopCountries)
		assert.Equal(t, []models.SpeciesMetric{{Species: "Acinetobacter baumannii", Count: 1}},
			result.SpeciesBreakdown)
	})

	t.Run("Success - Sorted Countries and Species", func(t *testing.T) {
		samples := []models.Sample{mockSample, mockSample}
		samples[1].Country.Code = "ARG"

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string, userID uuid.UUID, limit, offset int) ([]models.Sample, int64, error) {
				return samples, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID, filter models.AnalysisFilter, limit, offset int) ([]models.Analysis, int64, error) {
				return nil, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context, filter models.AdminUserFilter, limit, offset int) ([]models.User, int64, error) {
				return nil, 0, nil
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo, nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		assert.ElementsMatch(t, []models.CountryMetric{
			{Country: mockCountryCode, Count: 1},
			{Country: "ARG", Count: 1},
		}, result.TopCountries)
	})

	t.Run("Success - Aggregates Species and Genes", func(t *testing.T) {
		mockDone := testmodels.CreateMockAnalysis()
		mockDone.Status = models.AnalysisStatusDone
		mockDone.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","acquired_resistance":["blaOXA-23","armA"]}`)
		mockDuplicate := testmodels.CreateMockAnalysis()
		mockDuplicate.Status = models.AnalysisStatusDone
		mockDuplicate.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","acquired_resistance":["blaOXA-23"]}`)
		mockInvalid := testmodels.CreateMockAnalysis()
		mockInvalid.Status = models.AnalysisStatusDone
		mockInvalid.Metrics = []byte(`{"primary_species":`)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string, userID uuid.UUID, limit, offset int) ([]models.Sample, int64, error) {
				return []models.Sample{mockSample}, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID, filter models.AnalysisFilter, limit, offset int) ([]models.Analysis, int64, error) {
				return []models.Analysis{mockDone, mockDuplicate, mockInvalid}, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context, filter models.AdminUserFilter, limit, offset int) ([]models.User, int64, error) {
				return []models.User{}, 0, nil
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo, nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.TotalSpecies)
		assert.Equal(t, int64(2), result.TotalResistance)
		assert.Equal(t, int64(3), result.TotalAnalyses)
		assert.Equal(t, int64(3), result.AnalysesByStatus.Done)
	})

	t.Run("Non-Done Analyses Excluded", func(t *testing.T) {
		mockDone := testmodels.CreateMockAnalysis()
		mockDone.Status = models.AnalysisStatusDone
		mockDone.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","acquired_resistance":["blaOXA-23"]}`)
		mockPending := testmodels.CreateMockAnalysis()
		mockPending.Status = models.AnalysisStatusPending
		mockPending.Metrics = []byte(`{"primary_species":"Other","acquired_resistance":["x"]}`)
		mockEmptyResult := testmodels.CreateMockAnalysis()
		mockEmptyResult.Status = models.AnalysisStatusDone
		mockEmptyResult.Metrics = []byte(`{}`)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string, userID uuid.UUID, limit, offset int) ([]models.Sample, int64, error) {
				return nil, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID, filter models.AnalysisFilter, limit, offset int) ([]models.Analysis, int64, error) {
				return []models.Analysis{mockDone, mockPending, mockEmptyResult}, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context, filter models.AdminUserFilter, limit, offset int) ([]models.User, int64, error) {
				return nil, 0, nil
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo, nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.TotalSpecies)
		assert.Equal(t, int64(1), result.TotalResistance)
		assert.Equal(t, int64(2), result.AnalysesByStatus.Done)
		assert.Equal(t, int64(3), result.TotalAnalyses)
	})

	t.Run("Analysis Type Metrics", func(t *testing.T) {
		t0 := time.Date(2024, time.May, 11, 0, 0, 0, 0, time.UTC)
		at := func(sec int) *time.Time {
			ts := t0.Add(time.Duration(sec) * time.Second)
			return &ts
		}
		newAnalysis := func(analysisType models.AnalysisType,
			status models.AnalysisStatus, startSec, finishSec int) models.Analysis {
			a := testmodels.CreateMockAnalysis()
			a.Type = analysisType
			a.Status = status
			a.StartedAt = at(startSec)
			if finishSec > startSec {
				a.FinishedAt = at(finishSec)
			}
			return a
		}

		fastqcDone1 := newAnalysis(models.AnalysisTypeFastQC,
			models.AnalysisStatusDone, 0, 100)
		fastqcDone2 := newAnalysis(models.AnalysisTypeFastQC,
			models.AnalysisStatusDone, 0, 300)
		fastqcFailed := newAnalysis(models.AnalysisTypeFastQC,
			models.AnalysisStatusFailed, 0, 5)
		genomeDone := newAnalysis(models.AnalysisTypeGenome,
			models.AnalysisStatusDone, 0, 50)
		completePending := newAnalysis(models.AnalysisTypeComplete,
			models.AnalysisStatusPending, 0, 0)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID, limit, offset int) ([]models.Sample,
				int64, error) {
				return nil, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID,
				filter models.AnalysisFilter, limit, offset int) (
				[]models.Analysis, int64, error) {
				return []models.Analysis{fastqcDone1, fastqcDone2,
					fastqcFailed, genomeDone, completePending}, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter, limit, offset int) (
				[]models.User, int64, error) {
				return nil, 0, nil
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		if !assert.Len(t, result.AnalysesByType, 3) {
			return
		}
		assert.Equal(t, models.AnalysisTypeMetric{
			Type: "FASTQC", Count: 3,
			AvgDurationSeconds: 200, P95DurationSeconds: 300,
		}, result.AnalysesByType[0])
		assert.Equal(t, models.AnalysisTypeMetric{
			Type: "GENOME", Count: 1,
			AvgDurationSeconds: 50, P95DurationSeconds: 50,
		}, result.AnalysesByType[1])
		assert.Equal(t, models.AnalysisTypeMetric{
			Type: "COMPLETE", Count: 1,
		}, result.AnalysesByType[2])
	})

	t.Run("Success - Network counts", func(t *testing.T) {
		sampleIn := testmodels.CreateMockSample()
		sampleIn.InNetwork = true
		sampleOut := testmodels.CreateMockSample()
		sampleOut.InNetwork = false
		userIn := testmodels.NewLoginUser()
		userIn.IsPartOfNetwork = true
		userOut := testmodels.NewLoginUser()
		userOut.IsPartOfNetwork = false

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID, limit, offset int) ([]models.Sample,
				int64, error) {
				return []models.Sample{sampleIn, sampleOut, sampleIn}, 0,
					nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID,
				filter models.AnalysisFilter, limit, offset int) (
				[]models.Analysis, int64, error) {
				return nil, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter, limit, offset int) (
				[]models.User, int64, error) {
				return []models.User{userIn, userOut}, 0, nil
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		assert.Equal(t, models.NetworkMetric{InNetwork: 2, OutNetwork: 1},
			result.SamplesByNetwork)
		assert.Equal(t, models.NetworkMetric{InNetwork: 1, OutNetwork: 1},
			result.UsersByNetwork)
	})

	t.Run("Success - Top Resistance Genes", func(t *testing.T) {
		mockDone := testmodels.CreateMockAnalysis()
		mockDone.Status = models.AnalysisStatusDone
		mockDone.Metrics = []byte(`{"primary_species":"Klebsiella pneumoniae","acquired_resistance":["blaOXA-23","armA"]}`)
		mockDuplicate := testmodels.CreateMockAnalysis()
		mockDuplicate.Status = models.AnalysisStatusDone
		mockDuplicate.Metrics = []byte(`{"acquired_resistance":["blaOXA-23"]}`)
		mockInvalid := testmodels.CreateMockAnalysis()
		mockInvalid.Status = models.AnalysisStatusDone
		mockInvalid.Metrics = []byte(`{"acquired_resistance":`)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID, limit, offset int) ([]models.Sample,
				int64, error) {
				return nil, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID,
				filter models.AnalysisFilter, limit, offset int) (
				[]models.Analysis, int64, error) {
				return []models.Analysis{mockDone, mockDuplicate,
					mockInvalid}, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter, limit, offset int) (
				[]models.User, int64, error) {
				return nil, 0, nil
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []models.GeneMetric{
			{Gene: "blaOXA-23", Count: 2},
			{Gene: "armA", Count: 1},
		}, result.TopResistanceGenes)
		assert.Equal(t, int64(2), result.TotalResistance)
	})

	t.Run("Success - Samples by Month", func(t *testing.T) {
		now := time.Now()
		mockNow := testmodels.CreateMockSample()
		mockNow.CreatedAt = now
		mockPast := testmodels.CreateMockSample()
		mockPast.CreatedAt = now.AddDate(0, -2, 0)
		mockOld := testmodels.CreateMockSample()
		mockOld.CreatedAt = now.AddDate(0, -13, 0)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID, limit, offset int) ([]models.Sample,
				int64, error) {
				return []models.Sample{mockNow, mockPast, mockOld}, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID,
				filter models.AnalysisFilter, limit, offset int) (
				[]models.Analysis, int64, error) {
				return nil, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter, limit, offset int) (
				[]models.User, int64, error) {
				return nil, 0, nil
			},
		}
		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		if assert.Len(t, result.SamplesByMonth, 12) {
			for i := 1; i < len(result.SamplesByMonth); i++ {
				assert.True(t,
					result.SamplesByMonth[i].Month > result.SamplesByMonth[i-1].Month,
					"months must be ascending")
			}
		}
		byMonth := make(map[string]models.MonthMetric, 12)
		for _, m := range result.SamplesByMonth {
			byMonth[m.Month] = m
		}
		assert.Equal(t, int64(1),
			byMonth[now.Format("2006-01")].Count)
		assert.Equal(t, int64(1),
			byMonth[now.AddDate(0, -2, 0).Format("2006-01")].Count)
		assert.Equal(t, int64(0),
			byMonth[now.AddDate(0, -13, 0).Format("2006-01")].Count)
	})

	t.Run("Success - Queue Wait", func(t *testing.T) {
		t0 := time.Date(2024, time.May, 11, 0, 0, 0, 0, time.UTC)
		newAnalysis := func(waitSec int, started bool) models.Analysis {
			a := testmodels.CreateMockAnalysis()
			a.CreatedAt = t0
			a.StartedAt = nil
			if started {
				startedAt := t0.Add(time.Duration(waitSec) * time.Second)
				a.StartedAt = &startedAt
			}
			return a
		}

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID, limit, offset int) ([]models.Sample,
				int64, error) {
				return nil, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, userID uuid.UUID,
				filter models.AnalysisFilter, limit, offset int) (
				[]models.Analysis, int64, error) {
				return []models.Analysis{
					newAnalysis(100, true),
					newAnalysis(300, true),
					newAnalysis(0, false),
				}, 0, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter, limit, offset int) (
				[]models.User, int64, error) {
				return nil, 0, nil
			},
		}
		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			nil)
		result, err := svc.GetMetrics(ctx)

		assert.NoError(t, err)
		assert.Equal(t, models.QueueWaitMetric{
			AvgSeconds: 200,
			P95Seconds: 300,
		}, result.QueueWait)
	})

	t.Run("Error", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string, userID uuid.UUID, limit, offset int) ([]models.Sample, int64, error) {
				return nil, 0, services.ErrInternal
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{}
		userRepo := &mocks.MockUserRepository{}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			mockLogger)
		result, err := svc.GetMetrics(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestMetricsGetUserMetrics(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	buildRepos := func(samples []models.Sample, analyses []models.Analysis) (
		*mocks.MockSampleRepository, *mocks.MockAnalysisRepository,
		*mocks.MockUserRepository) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				repoUserID uuid.UUID, limit, offset int) ([]models.Sample,
				int64, error) {
				assert.Equal(t, userID, repoUserID)
				return samples, 0, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, repoUserID uuid.UUID,
				filter models.AnalysisFilter, limit, offset int) (
				[]models.Analysis, int64, error) {
				assert.Equal(t, userID, repoUserID)
				return analyses, 0, nil
			},
		}
		return sampleRepo, analysisRepo, &mocks.MockUserRepository{}
	}

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		mockDone := testmodels.CreateMockAnalysis()
		mockDone.Type = models.AnalysisTypeFastQC
		mockDone.Status = models.AnalysisStatusDone
		mockDone.UserID = userID
		mockDone.Metrics = []byte(`{"primary_species":"Acinetobacter sp","acquired_resistance":["g01","g01","g01","g02","g02","g03","g04","g05","g06","g07","g08","g09","g10","g11"]}`)
		mockSampleNow := testmodels.CreateMockSample()
		mockSampleNow.UserID = userID
		mockSampleNow.CreatedAt = now
		mockSamplePast := testmodels.CreateMockSample()
		mockSamplePast.UserID = userID
		mockSamplePast.CreatedAt = now.AddDate(0, -2, 0)

		sampleRepo, analysisRepo, userRepo := buildRepos(
			[]models.Sample{mockSampleNow, mockSamplePast},
			[]models.Analysis{mockDone})

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			nil)
		result, err := svc.GetUserMetrics(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, models.AnalysesByStatus{Done: 1},
			result.AnalysesByStatus)
		if assert.Len(t, result.AnalysesByType, 3) {
			assert.Equal(t, models.AnalysisTypeMetric{
				Type: string(models.AnalysisTypeFastQC), Count: 1},
				result.AnalysesByType[0])
		}
		assert.Equal(t, []models.SpeciesMetric{
			{Species: "Acinetobacter sp", Count: 1},
		}, result.TopSpecies)
		if assert.Len(t, result.TopResistanceGenes, 11) {
			assert.Equal(t, models.GeneMetric{Gene: "g01", Count: 3},
				result.TopResistanceGenes[0])
			assert.Equal(t, models.GeneMetric{Gene: "g02", Count: 2},
				result.TopResistanceGenes[1])
		}
		if assert.Len(t, result.SamplesByMonth, 12) {
			byMonth := make(map[string]models.MonthMetric, 12)
			for _, m := range result.SamplesByMonth {
				byMonth[m.Month] = m
			}
			assert.Equal(t, int64(1),
				byMonth[now.Format("2006-01")].Count)
			assert.Equal(t, int64(1),
				byMonth[now.AddDate(0, -2, 0).Format("2006-01")].Count)
		}
	})

	t.Run("Success - Empty", func(t *testing.T) {
		sampleRepo, analysisRepo, userRepo := buildRepos(nil, nil)

		svc := services.NewMetricsService(sampleRepo, analysisRepo, userRepo,
			nil)
		result, err := svc.GetUserMetrics(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, models.AnalysesByStatus{},
			result.AnalysesByStatus)
		if assert.Len(t, result.AnalysesByType, 3) {
			for _, entry := range result.AnalysesByType {
				assert.Equal(t, int64(0), entry.Count)
			}
		}
		assert.Equal(t, []models.SpeciesMetric{}, result.TopSpecies)
		assert.Equal(t, []models.GeneMetric{},
			result.TopResistanceGenes)
		if assert.Len(t, result.SamplesByMonth, 12) {
			for _, m := range result.SamplesByMonth {
				assert.Equal(t, int64(0), m.Count)
			}
		}
	})

	t.Run("Error - Analysis Repo", func(t *testing.T) {
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		sampleRepo := &mocks.MockSampleRepository{}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context, repoUserID uuid.UUID,
				filter models.AnalysisFilter, limit, offset int) (
				[]models.Analysis, int64, error) {
				return nil, 0, gorm.ErrInvalidTransaction
			},
		}

		svc := services.NewMetricsService(sampleRepo, analysisRepo,
			&mocks.MockUserRepository{}, mockLogger)
		result, err := svc.GetUserMetrics(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sample Repo", func(t *testing.T) {
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				repoUserID uuid.UUID, limit, offset int) ([]models.Sample,
				int64, error) {
				return nil, 0, gorm.ErrInvalidTransaction
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{}

		svc := services.NewMetricsService(sampleRepo, analysisRepo,
			&mocks.MockUserRepository{}, mockLogger)
		result, err := svc.GetUserMetrics(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}
