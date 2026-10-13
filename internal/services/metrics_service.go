package services

import (
	"context"
	"encoding/json"
	"slices"
	"sort"
	"time"

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
	samples, _, err := s.SampleRepo.GetSamples(ctx, "", uuid.Nil, 0, 0)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"MetricsService", "GetMetrics", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	analyses, _, err := s.AnalysisRepo.GetAnalyses(ctx, uuid.Nil,
		models.AnalysisFilter{}, 0, 0)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
			"MetricsService", "GetMetrics", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	users, _, err := s.UserRepo.GetUsers(ctx, models.AdminUserFilter{}, 0, 0)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(ctx,
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
		UsersByNetwork:   countUsersByNetwork(users),
		SamplesByNetwork: countSamplesByNetwork(samples),
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

	metrics.AnalysesByType = analysisTypeMetrics(analyses)

	species, genes := uniqueSpeciesResults(analyses)
	metrics.TotalSpecies = int64(len(species))
	metrics.TotalResistance = int64(len(genes))
	metrics.TopResistanceGenes = sortedGeneMetrics(genes)
	metrics.SamplesByMonth = samplesByMonth(samples)
	metrics.QueueWait = queueWait(analyses)

	return metrics, nil
}

func countUsersByNetwork(users []models.User) models.NetworkMetric {
	var metric models.NetworkMetric
	for _, user := range users {
		if user.IsPartOfNetwork {
			metric.InNetwork++
		} else {
			metric.OutNetwork++
		}
	}

	return metric
}

func countSamplesByNetwork(samples []models.Sample) models.NetworkMetric {
	var metric models.NetworkMetric
	for _, sample := range samples {
		if sample.InNetwork {
			metric.InNetwork++
		} else {
			metric.OutNetwork++
		}
	}

	return metric
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

func uniqueSpeciesResults(analyses []models.Analysis) (map[string]struct{}, map[string]int64) {
	species := make(map[string]struct{})
	genes := make(map[string]int64)
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
				genes[gene]++
			}
		}
	}
	return species, genes
}

func sortedGeneMetrics(counts map[string]int64) []models.GeneMetric {
	metrics := make([]models.GeneMetric, 0, len(counts))
	for gene, count := range counts {
		metrics = append(metrics, models.GeneMetric{
			Gene: gene, Count: count,
		})
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Count > metrics[j].Count
	})
	return metrics
}

func samplesByMonth(samples []models.Sample) []models.MonthMetric {
	start := time.Now().AddDate(0, -11, 0)

	counts := make(map[string]int64)
	for _, sample := range samples {
		counts[sample.CreatedAt.Format("2006-01")]++
	}

	metrics := make([]models.MonthMetric, 0, 12)
	for i := range 12 {
		key := start.AddDate(0, i, 0).Format("2006-01")
		metrics = append(metrics, models.MonthMetric{
			Month: key, Count: counts[key],
		})
	}
	return metrics
}

func queueWait(analyses []models.Analysis) models.QueueWaitMetric {
	var waits []time.Duration
	for _, analysis := range analyses {
		if analysis.StartedAt != nil {
			waits = append(waits, analysis.StartedAt.Sub(analysis.CreatedAt))
		}
	}

	return models.QueueWaitMetric{
		AvgSeconds: avgDuration(waits),
		P95Seconds: p95Duration(waits),
	}
}

func analysisTypeMetrics(analyses []models.Analysis) []models.AnalysisTypeMetric {
	type bucket struct {
		count     int64
		durations []time.Duration
	}
	buckets := make(map[models.AnalysisType]*bucket, len(models.AnalysisTypes))
	for _, analysis := range analyses {
		b := buckets[analysis.Type]
		if b == nil {
			b = &bucket{}
			buckets[analysis.Type] = b
		}
		b.count++
		if analysis.Status == models.AnalysisStatusDone &&
			analysis.StartedAt != nil && analysis.FinishedAt != nil {
			b.durations = append(b.durations,
				analysis.FinishedAt.Sub(*analysis.StartedAt))
		}
	}
	metrics := make([]models.AnalysisTypeMetric, 0, len(models.AnalysisTypes))
	for _, analysisType := range models.AnalysisTypes {
		b := buckets[analysisType]
		if b == nil {
			metrics = append(metrics, models.AnalysisTypeMetric{
				Type: string(analysisType),
			})
			continue
		}
		metrics = append(metrics, models.AnalysisTypeMetric{
			Type:               string(analysisType),
			Count:              b.count,
			AvgDurationSeconds: avgDuration(b.durations),
			P95DurationSeconds: p95Duration(b.durations),
		})
	}
	return metrics
}

func avgDuration(durations []time.Duration) float64 {
	if len(durations) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	return sum.Seconds() / float64(len(durations))
}

func p95Duration(durations []time.Duration) float64 {
	if len(durations) == 0 {
		return 0
	}
	sorted := append([]time.Duration(nil), durations...)
	slices.Sort(sorted)
	idx := (len(sorted)*95+99)/100 - 1
	return sorted[idx].Seconds()
}
