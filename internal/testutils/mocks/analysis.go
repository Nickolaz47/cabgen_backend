package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockAnalysisRepository struct {
	GetAnalysesFunc func(ctx context.Context, userID uuid.UUID) (
		[]models.Analysis, error)
	GetAnalysisByIDFunc func(ctx context.Context, analysisID uuid.UUID) (
		*models.Analysis, error)
	CreateAnalysisFunc func(ctx context.Context,
		analysis *models.Analysis) error
	UpdateAnalysisFunc func(ctx context.Context,
		analysis *models.Analysis) error
	DeleteAnalysisFunc func(ctx context.Context,
		analysis *models.Analysis) error
}

func (r *MockAnalysisRepository) GetAnalyses(ctx context.Context,
	userID uuid.UUID) ([]models.Analysis, error) {
	if r.GetAnalysesFunc != nil {
		return r.GetAnalysesFunc(ctx, userID)
	}
	return nil, nil
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

func (r *MockAnalysisRepository) DeleteAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	if r.DeleteAnalysisFunc != nil {
		return r.DeleteAnalysisFunc(ctx, analysis)
	}

	return nil
}

type MockAnalysisService struct {
	FindAllFunc func(ctx context.Context, input string, userID uuid.UUID) (
		[]models.AnalysisResponse, error)
	FindByIDFunc func(ctx context.Context, analysisID, userID uuid.UUID) (
		*models.AnalysisResponse, error)
	CreateFunc func(ctx context.Context, input models.AnalysisCreateDTO) (
		*models.AnalysisResponse, error)
	DeleteFunc func(ctx context.Context, analysisID, userID uuid.UUID) error
}

func (s *MockAnalysisService) FindAll(ctx context.Context, input string,
	userID uuid.UUID) (
	[]models.AnalysisResponse, error) {
	if s.FindAllFunc != nil {
		return s.FindAllFunc(ctx, input, userID)
	}

	return nil, nil
}

func (s *MockAnalysisService) FindByID(ctx context.Context, analysisID,
	userID uuid.UUID) (
	*models.AnalysisResponse, error) {
	if s.FindByIDFunc != nil {
		return s.FindByIDFunc(ctx, analysisID, userID)
	}

	return nil, nil
}

func (s *MockAnalysisService) Create(ctx context.Context,
	input models.AnalysisCreateDTO) (
	*models.AnalysisResponse, error) {
	if s.CreateFunc != nil {
		return s.CreateFunc(ctx, input)
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
