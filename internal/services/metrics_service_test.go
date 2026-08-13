package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMetricsGetMetrics(t *testing.T) {
	ctx := context.Background()

	mockSample := testmodels.CreateMockSample()
	mockUser := testmodels.NewLoginUser()

	var mockCountryCode = "BRA"

	t.Run("Success", func(t *testing.T) {
		mockDone := testmodels.CreateMockAnalysis()
		mockDone.Status = models.AnalysisStatusDone
		mockDone.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","gene":["blaOXA-23"]}`)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID) ([]models.Sample, error) {
				return []models.Sample{mockSample}, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return []models.Analysis{mockDone}, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{mockUser}, nil
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
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID) ([]models.Sample, error) {
				return samples, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return nil, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return nil, nil
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
		mockDone.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","gene":["blaOXA-23","armA"]}`)
		mockDuplicate := testmodels.CreateMockAnalysis()
		mockDuplicate.Status = models.AnalysisStatusDone
		mockDuplicate.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","gene":["blaOXA-23"]}`)
		mockInvalid := testmodels.CreateMockAnalysis()
		mockInvalid.Status = models.AnalysisStatusDone
		mockInvalid.Metrics = []byte(`{"primary_species":`)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID) ([]models.Sample, error) {
				return []models.Sample{mockSample}, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return []models.Analysis{mockDone, mockDuplicate, mockInvalid}, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{}, nil
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
		mockDone.Metrics = []byte(`{"primary_species":"Acinetobacter baumannii","gene":["blaOXA-23"]}`)
		mockPending := testmodels.CreateMockAnalysis()
		mockPending.Status = models.AnalysisStatusPending
		mockPending.Metrics = []byte(`{"primary_species":"Other","gene":["x"]}`)
		mockEmptyResult := testmodels.CreateMockAnalysis()
		mockEmptyResult.Status = models.AnalysisStatusDone
		mockEmptyResult.Metrics = []byte(`{}`)

		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID) ([]models.Sample, error) {
				return nil, nil
			},
		}
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return []models.Analysis{mockDone, mockPending, mockEmptyResult}, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return nil, nil
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

	t.Run("Error", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID) ([]models.Sample, error) {
				return nil, services.ErrInternal
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
