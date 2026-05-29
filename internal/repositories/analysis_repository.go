package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnalysisRepository interface {
	GetAnalyses(ctx context.Context, userID uuid.UUID) (
		[]models.Analysis, error)
	GetAnalysesByIDs(ctx context.Context, analysisIDs []uuid.UUID,
		userID uuid.UUID) ([]models.Analysis, error)
	GetAnalysisByID(ctx context.Context, analysisID uuid.UUID) (
		*models.Analysis, error)
	CreateAnalysis(ctx context.Context, analysis *models.Analysis) error
	UpdateAnalysis(ctx context.Context, analysis *models.Analysis) error
	DeleteAnalysis(ctx context.Context, analysis *models.Analysis) error
}

type analysisRepo struct {
	DB *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) AnalysisRepository {
	return &analysisRepo{
		DB: db,
	}
}

func (r *analysisRepo) GetAnalyses(ctx context.Context, userID uuid.UUID) (
	[]models.Analysis, error) {
	var analyses []models.Analysis

	query := r.DB.WithContext(ctx).Preload("Sample").Preload("User")
	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&analyses).Error; err != nil {
		return nil, err
	}

	return analyses, nil
}

func (r *analysisRepo) GetAnalysesByIDs(ctx context.Context,
	analysisIDs []uuid.UUID, userID uuid.UUID) (
	[]models.Analysis, error) {
	var analyses []models.Analysis

	query := r.DB.WithContext(ctx).Preload("Sample").Preload("User").
		Where("id in ?", analysisIDs)

	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&analyses).Error; err != nil {
		return nil, err
	}

	return analyses, nil
}

func (r *analysisRepo) GetAnalysisByID(ctx context.Context,
	analysisID uuid.UUID) (*models.Analysis, error) {
	var analysis models.Analysis
	if err := r.DB.WithContext(ctx).Preload("Sample").Preload("User").
	Where("id = ?", analysisID).First(
		&analysis).Error; err != nil {
		return nil, err
	}

	return &analysis, nil
}

func (r *analysisRepo) CreateAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	return r.DB.WithContext(ctx).Create(analysis).Error
}

func (r *analysisRepo) UpdateAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	return r.DB.WithContext(ctx).Save(analysis).Error
}

func (r *analysisRepo) DeleteAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	return r.DB.WithContext(ctx).Delete(analysis).Error
}
