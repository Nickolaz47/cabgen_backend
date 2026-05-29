package services

import (
	"context"
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AnalysisService interface {
	FindAll(ctx context.Context, userID uuid.UUID) (
		[]models.AnalysisResponse, error)
	FindByID(ctx context.Context, analysisID, userID uuid.UUID) (
		*models.AnalysisResponse, error)
	Create(ctx context.Context, input models.AnalysisCreateDTO) (
		*models.AnalysisResponse, error)
	Delete(ctx context.Context, analysisID, userID uuid.UUID) error
}

type analysisService struct {
	Repo       repositories.AnalysisRepository
	SampleRepo repositories.SampleRepository
	UserRepo   repositories.UserRepository
	Logger     *zap.Logger
}

func NewAnalysisService(
	repo repositories.AnalysisRepository,
	sampleRepo repositories.SampleRepository,
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) AnalysisService {
	return &analysisService{
		Repo:       repo,
		SampleRepo: sampleRepo,
		UserRepo:   userRepo,
		Logger:     logger,
	}
}

func (s *analysisService) FindAll(ctx context.Context, userID uuid.UUID) (
	[]models.AnalysisResponse, error) {
	analyses, err := s.Repo.GetAnalyses(ctx, userID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "FindAll",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	responses := make([]models.AnalysisResponse, len(analyses))
	for i, analysis := range analyses {
		responses[i] = analysis.ToResponse()
	}

	return responses, nil
}

func (s *analysisService) FindByID(ctx context.Context, analysisID,
	userID uuid.UUID) (
	*models.AnalysisResponse, error) {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "FindByID", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "FindByID",
			logging.DatabaseError, err)...)
		return nil, ErrInternal
	}

	if userID != uuid.Nil && userID != analysis.UserID {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "FindByID", logging.Unauthorized, err,
		)...)
		return nil, ErrUnauthorized
	}

	response := analysis.ToResponse()
	return &response, nil
}

func (s *analysisService) Create(ctx context.Context,
	input models.AnalysisCreateDTO) (
	*models.AnalysisResponse, error) {
	sample, err := s.SampleRepo.GetSampleByID(ctx, input.SampleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AnalysisService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrSampleNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AnalysisService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	if input.Type == models.AnalysisTypeFastQC &&
		sample.Fastq1 == nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AnalysisService", "Create",
				logging.MissingFileError, ErrMissingFastq1,
			)...)
		return nil, ErrMissingFastq1
	} else if input.Type == models.AnalysisTypeFastQC &&
		sample.Fastq2 == nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AnalysisService", "Create",
				logging.MissingFileError, ErrMissingFastq2,
			)...)
		return nil, ErrMissingFastq2
	}

	if input.Type == models.AnalysisTypeComplete && (
		sample.Fastq1 == nil || sample.Fastq2 == nil) &&
		sample.Fasta != nil {
		input.Type = models.AnalysisTypeGenome
	}

	user, err := s.UserRepo.GetUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AnalysisService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrUserNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AnalysisService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	analysis := models.Analysis{
		Type:     input.Type,
		Status:   models.AnalysisStatusPending,
		SampleID: sample.ID,
		UserID:   input.UserID,
	}

	analysis.Sample = *sample
	analysis.User = *user

	if err := s.Repo.CreateAnalysis(ctx, &analysis); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AnalysisService", "Create",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := analysis.ToResponse()
	return &response, nil
}

func (s *analysisService) Delete(ctx context.Context,
	analysisID, userID uuid.UUID) error {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Delete", logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if userID != uuid.Nil && userID != analysis.UserID {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Delete", logging.Unauthorized, err,
		)...)
		return ErrUnauthorized
	}

	if analysis.Status == models.AnalysisStatusRunning {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrDeleteRunningAnalysis
	}

	if err := s.Repo.DeleteAnalysis(ctx, analysis); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}
