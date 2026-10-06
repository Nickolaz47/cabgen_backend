package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnalysisRepository interface {
	GetAnalyses(ctx context.Context, userID uuid.UUID,
		filter models.AnalysisFilter) ([]models.Analysis, error)
	GetAnalysesByIDs(ctx context.Context, analysisIDs []uuid.UUID,
		userID uuid.UUID) ([]models.Analysis, error)
	GetDashboardAnalyses(ctx context.Context) ([]models.Analysis, error)
	GetAnalysisByID(ctx context.Context, analysisID uuid.UUID) (
		*models.Analysis, error)
	CreateAnalysis(ctx context.Context, analysis *models.Analysis) error
	UpdateAnalysis(ctx context.Context, analysis *models.Analysis) error
	UpdateSample(ctx context.Context, sample *models.Sample) error
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

func (r *analysisRepo) GetAnalyses(ctx context.Context, userID uuid.UUID,
	filter models.AnalysisFilter) (
	[]models.Analysis, error) {
	var analyses []models.Analysis

	query := r.DB.WithContext(ctx).Preload("Sample").Preload("User")
	// Collaborator path
	if userID != uuid.Nil {
		query = query.Where("analyses.user_id = ?", userID)

		if filter.OriginCode != "" {
			like := "%" + filter.OriginCode + "%"
			query = query.Joins("JOIN samples ON samples.id"+
				" = analyses.sample_id").Where(
				"LOWER(samples.origin_code) LIKE LOWER(?)", like,
			)
		}

		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}
	}

	// Admin path
	if userID == uuid.Nil {
		if filter.OriginCode != "" {
			like := "%" + filter.OriginCode + "%"
			query = query.Joins("JOIN samples ON samples.id"+
				" = analyses.sample_id").Where(
				"LOWER(samples.origin_code) LIKE LOWER(?)", like,
			)
		}

		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}

		if filter.Username != "" {
			query = query.Joins("JOIN users ON users.id"+
				" = analyses.user_id").Where("users.username = ?",
				filter.Username)
		}
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

func (r *analysisRepo) GetDashboardAnalyses(ctx context.Context) (
	[]models.Analysis, error) {
	var analyses []models.Analysis

	query := `analyses.id IN (
		SELECT a2.id FROM analyses a2
		JOIN samples s2 ON s2.id = a2.sample_id
		WHERE s2.in_network = true
			AND a2.status = 'DONE'
			AND a2.type IN ('GENOME', 'COMPLETE')
			AND NOT EXISTS (
				SELECT 1 FROM analyses a3
				JOIN samples s3 ON s3.id = a3.sample_id
				WHERE s3.in_network = true
					AND a3.status = 'DONE'
					AND a3.type IN ('GENOME', 'COMPLETE')
					AND a3.sample_id = a2.sample_id
					AND (a3.created_at > a2.created_at
						OR (a3.created_at = a2.created_at
							AND a3.id > a2.id))
			)
	)`

	if err := r.DB.WithContext(ctx).
		Preload("Sample.SampleSource").
		Where(query).
		Order("analyses.created_at DESC").
		Find(&analyses).Error; err != nil {
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

func (r *analysisRepo) UpdateSample(ctx context.Context,
	sample *models.Sample) error {
	return r.DB.WithContext(ctx).Save(sample).Error
}

func (r *analysisRepo) DeleteAnalysis(ctx context.Context,
	analysis *models.Analysis) error {
	return r.DB.WithContext(ctx).Delete(analysis).Error
}
