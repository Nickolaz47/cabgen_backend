package services

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MetricsService interface {
	GetMetrics(ctx context.Context) (*models.AdminMetricsResponse, error)
}

type metricsService struct {
	SampleRepo   repositories.SampleRepository
	AnalysisRepo repositories.AnalysisRepository
	UserRepo     repositories.UserRepository
	Logger       *zap.Logger
}

func NewMetricsService(
	sampleRepo repositories.SampleRepository,
	analysisRepo repositories.AnalysisRepository,
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) MetricsService {
	return &metricsService{
		SampleRepo:   sampleRepo,
		AnalysisRepo: analysisRepo,
		UserRepo:     userRepo,
		Logger:       logger,
	}
}

func (s *metricsService) GetMetrics(ctx context.Context) (
	*models.AdminMetricsResponse, error) {
	samples, err := s.SampleRepo.GetSamples(ctx, "", uuid.Nil)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MetricsService", "GetMetrics", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	analyses, err := s.AnalysisRepo.GetAnalyses(ctx, uuid.Nil)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MetricsService", "GetMetrics", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	users, err := s.UserRepo.GetUsers(ctx, models.AdminUserFilter{})
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"MetricsService", "GetMetrics", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	metrics := &models.AdminMetricsResponse{
		PublicMetricsResponse: models.PublicMetricsResponse{
			TotalSamples:   int64(len(samples)),
			TotalCountries: int64(len(uniqueCountryCodes(samples))),
		},
		TotalUsers:       int64(len(users)),
		TotalAnalyses:    int64(len(analyses)),
		TopCountries:     countryMetrics(samples),
		SpeciesBreakdown: speciesMetrics(analyses),
	}

	for _, analysis := range analyses {
		switch analysis.Status {
		case models.AnalysisStatusDone:
			metrics.AnalysesByStatus.Done++
		case models.AnalysisStatusRunning:
			metrics.AnalysesByStatus.Running++
		case models.AnalysisStatusPending:
			metrics.AnalysesByStatus.Pending++
		case models.AnalysisStatusFailed:
			metrics.AnalysesByStatus.Failed++
		}
	}

	species, genes := uniqueSpeciesResults(analyses)
	metrics.TotalSpecies = int64(len(species))
	metrics.TotalResistance = int64(len(genes))

	return metrics, nil
}

func uniqueCountryCodes(samples []models.Sample) map[string]struct{} {
	codes := make(map[string]struct{})
	for _, sample := range samples {
		if sample.Country.Code != "" {
			codes[sample.Country.Code] = struct{}{}
		}
	}
	return codes
}

func countryMetrics(samples []models.Sample) []models.CountryMetric {
	counts := make(map[string]int64)
	for _, sample := range samples {
		if sample.Country.Code != "" {
			counts[sample.Country.Code]++
		}
	}
	return sortedCountryMetrics(counts)
}

func sortedCountryMetrics(counts map[string]int64) []models.CountryMetric {
	metrics := make([]models.CountryMetric, 0, len(counts))
	for code, count := range counts {
		metrics = append(metrics, models.CountryMetric{
			Country: code, Count: count,
		})
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Count > metrics[j].Count
	})
	return metrics
}

func speciesMetrics(analyses []models.Analysis) []models.SpeciesMetric {
	counts := make(map[string]int64)
	for _, analysis := range analyses {
		if analysis.Status != models.AnalysisStatusDone || len(analysis.Metrics) == 0 {
			continue
		}
		var results models.AnalysisResults
		if err := json.Unmarshal(analysis.Metrics, &results); err != nil {
			continue
		}
		if results.PrimarySpeciesName != "" {
			counts[results.PrimarySpeciesName]++
		}
	}
	metrics := make([]models.SpeciesMetric, 0, len(counts))
	for species, count := range counts {
		metrics = append(metrics, models.SpeciesMetric{
			Species: species, Count: count,
		})
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Count > metrics[j].Count
	})
	return metrics
}

func uniqueSpeciesResults(analyses []models.Analysis) (map[string]struct{}, map[string]struct{}) {
	species := make(map[string]struct{})
	genes := make(map[string]struct{})
	for _, analysis := range analyses {
		if analysis.Status != models.AnalysisStatusDone || len(analysis.Metrics) == 0 {
			continue
		}
		var results models.AnalysisResults
		if err := json.Unmarshal(analysis.Metrics, &results); err != nil {
			continue
		}
		if results.PrimarySpeciesName != "" {
			species[results.PrimarySpeciesName] = struct{}{}
		}
		for _, gene := range results.AcquiredResistance {
			if gene != "" {
				genes[gene] = struct{}{}
			}
		}
	}
	return species, genes
}
