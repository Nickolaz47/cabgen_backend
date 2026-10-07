package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockAnalysisRepository struct {
	GetAnalysesFunc func(ctx context.Context, userID uuid.UUID,
		filter models.AnalysisFilter, limit, offset int) (
		[]models.Analysis, int64, error)
	GetAnalysesByIDsFunc func(ctx context.Context, analysisIDs []uuid.UUID,
		userID uuid.UUID) ([]models.Analysis, error)
	GetDashboardAnalysesFunc func(ctx context.Context) (
		[]models.Analysis, error)
	CountAnalysesBySampleFunc func(ctx context.Context,
		sampleID uuid.UUID) (int64, error)
	GetAnalysisByIDFunc func(ctx context.Context, analysisID uuid.UUID) (
		*models.Analysis, error)
	CreateAnalysisFunc func(ctx context.Context,
		analysis *models.Analysis) error
	UpdateAnalysisFunc func(ctx context.Context,
		analysis *models.Analysis) error
	UpdateSampleFunc func(ctx context.Context,
		sample *models.Sample) error
	DeleteAnalysisFunc func(ctx context.Context,
		analysis *models.Analysis) error
}

func (r *MockAnalysisRepository) GetAnalyses(ctx context.Context,
	userID uuid.UUID, filter models.AnalysisFilter, limit, offset int) (
	[]models.Analysis, int64, error) {
	if r.GetAnalysesFunc != nil {
		return r.GetAnalysesFunc(ctx, userID, filter, limit, offset)
	}
	return nil, 0, nil
}

func (r *MockAnalysisRepository) GetAnalysesByIDs(ctx context.Context,
	analysisIDs []uuid.UUID, userID uuid.UUID) ([]models.Analysis, error) {
	if r.GetAnalysesByIDsFunc != nil {
		return r.GetAnalysesByIDsFunc(ctx, analysisIDs, userID)
	}

	return nil, nil
}

func (r *MockAnalysisRepository) GetDashboardAnalyses(ctx context.Context) (
	[]models.Analysis, error) {
	if r.GetDashboardAnalysesFunc != nil {
		return r.GetDashboardAnalysesFunc(ctx)
	}

	return nil, nil
}

func (r *MockAnalysisRepository) CountAnalysesBySample(ctx context.Context,
	sampleID uuid.UUID) (int64, error) {
	if r.CountAnalysesBySampleFunc != nil {
		return r.CountAnalysesBySampleFunc(ctx, sampleID)
	}

	return 0, nil
}

func (r *MockAnalysisRepository) GetAnalysisByID(ctx context.Context,
	analysisID uuid.UUID) (*models.Analysis, error) {
	if r.GetAnalysisByIDFunc != nil {
		return r.GetAnalysisByIDFunc(ctx, analysisID)
	}

	return nil, nil
}

func (r *MockAnalysisRepository) CreateAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	if r.CreateAnalysisFunc != nil {
		return r.CreateAnalysisFunc(ctx, analysis)
	}

	return nil
}

func (r *MockAnalysisRepository) UpdateAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	if r.UpdateAnalysisFunc != nil {
		return r.UpdateAnalysisFunc(ctx, analysis)
	}

	return nil
}

func (r *MockAnalysisRepository) UpdateSample(ctx context.Context,
	sample *models.Sample) error {
	if r.UpdateSampleFunc != nil {
		return r.UpdateSampleFunc(ctx, sample)
	}

	return nil
}

func (r *MockAnalysisRepository) DeleteAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	if r.DeleteAnalysisFunc != nil {
		return r.DeleteAnalysisFunc(ctx, analysis)
	}

	return nil
}

type MockAnalysisService struct {
	FindAllFunc func(ctx context.Context, userID uuid.UUID,
		filter models.AnalysisFilter, offset int, language string) (
		[]models.AnalysisResponse, int, error)
	FindByIDFunc func(ctx context.Context, analysisID, userID uuid.UUID,
		language string) (*models.AnalysisResponse, error)
	FindManyByIDsFunc func(ctx context.Context, analysisIDs []uuid.UUID,
		userID uuid.UUID, language string) ([]models.AnalysisResponse, error)
	CreateFunc func(ctx context.Context, input models.AnalysisCreateDTO,
		language string) (*models.AnalysisResponse, error)
	UpdateFunc func(ctx context.Context, analysisID uuid.UUID,
		input models.AdminAnalysisUpdateInput, language string) (
		*models.AnalysisResponse, error)
	DeleteFunc      func(ctx context.Context, analysisID, userID uuid.UUID) error
	DownloadZipFunc func(ctx context.Context, analysisID,
		userID uuid.UUID) (string, error)
	DownloadBatchTSVFunc func(ctx context.Context, analysisIDs []uuid.UUID,
		userID uuid.UUID, language string) ([]models.AnalysisResponse, error)
	DownloadDashboardTSVFunc func(ctx context.Context) (
		[]models.Analysis, error)
}

func (s *MockAnalysisService) FindAll(ctx context.Context, userID uuid.UUID,
	filter models.AnalysisFilter, offset int, language string) (
	[]models.AnalysisResponse, int, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx, userID, filter, offset, language)
	}

	return nil, 0, nil
}

func (s *MockAnalysisService) FindManyByIDs(ctx context.Context,
	analysisIDs []uuid.UUID, userID uuid.UUID, language string) (
	[]models.AnalysisResponse, error) {
	if s.FindManyByIDsFunc != nil {
		return s.FindManyByIDsFunc(ctx, analysisIDs, userID, language)
	}

	return nil, nil
}

func (s *MockAnalysisService) FindByID(ctx context.Context, analysisID,
	userID uuid.UUID, language string) (
	*models.AnalysisResponse, error) {
	if s.FindByIDFunc != nil {
		return s.FindByIDFunc(ctx, analysisID, userID, language)
	}

	return nil, nil
}

func (s *MockAnalysisService) Create(ctx context.Context,
	input models.AnalysisCreateDTO, language string) (
	*models.AnalysisResponse, error) {
	if s.CreateFunc != nil {
		return s.CreateFunc(ctx, input, language)
	}

	return nil, nil
}

func (s *MockAnalysisService) Update(ctx context.Context, analysisID uuid.UUID,
	input models.AdminAnalysisUpdateInput, language string) (
	*models.AnalysisResponse, error) {
	if s.UpdateFunc != nil {
		return s.UpdateFunc(ctx, analysisID, input, language)
	}

	return nil, nil
}

func (s *MockAnalysisService) Delete(ctx context.Context, analysisID,
	userID uuid.UUID) error {
	if s.DeleteFunc != nil {
		return s.DeleteFunc(ctx, analysisID, userID)
	}

	return nil
}

func (s *MockAnalysisService) DownloadZip(ctx context.Context, analysisID,
	userID uuid.UUID) (string, error) {
	if s.DownloadZipFunc != nil {
		return s.DownloadZipFunc(ctx, analysisID, userID)
	}

	return "", nil
}

func (s *MockAnalysisService) DownloadBatchTSV(ctx context.Context,
	analysisIDs []uuid.UUID, userID uuid.UUID, language string) (
	[]models.AnalysisResponse, error) {
	if s.DownloadBatchTSVFunc != nil {
		return s.DownloadBatchTSVFunc(ctx, analysisIDs, userID, language)
	}

	return nil, nil
}

func (s *MockAnalysisService) DownloadDashboardTSV(ctx context.Context) (
	[]models.Analysis, error) {
	if s.DownloadDashboardTSVFunc != nil {
		return s.DownloadDashboardTSVFunc(ctx)
	}

	return nil, nil
}
