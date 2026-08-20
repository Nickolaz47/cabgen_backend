package services

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AnalysisService interface {
	FindAll(ctx context.Context, userID uuid.UUID, filter models.AnalysisFilter,
		language string) (
		[]models.AnalysisResponse, error)
	FindByID(ctx context.Context, analysisID, userID uuid.UUID,
		language string) (*models.AnalysisResponse, error)
	FindManyByIDs(ctx context.Context, analysisIDs []uuid.UUID,
		userID uuid.UUID, language string) ([]models.AnalysisResponse, error)
	Create(ctx context.Context, input models.AnalysisCreateDTO,
		language string) (*models.AnalysisResponse, error)
	Update(ctx context.Context, analysisID uuid.UUID,
		input models.AdminAnalysisUpdateInput, language string) (
		*models.AnalysisResponse, error)
	Delete(ctx context.Context, analysisID, userID uuid.UUID) error
	DownloadZip(ctx context.Context, analysisID, userID uuid.UUID) (string,
		error)
	DownloadBatchTSV(ctx context.Context, analysisIDs []uuid.UUID,
		userID uuid.UUID, language string) ([]models.AnalysisResponse, error)
}

type analysisService struct {
	Repo        repositories.AnalysisRepository
	SampleRepo  repositories.SampleRepository
	UserRepo    repositories.UserRepository
	AsynqClient TaskEnqueuer
	Canceller   TaskCanceller
	Logger      *zap.Logger
	RootDir     string
}

func NewAnalysisService(
	repo repositories.AnalysisRepository,
	sampleRepo repositories.SampleRepository,
	userRepo repositories.UserRepository,
	asynqClient TaskEnqueuer,
	canceller TaskCanceller,
	logger *zap.Logger,
	rootDir string,
) AnalysisService {
	return &analysisService{
		Repo:        repo,
		SampleRepo:  sampleRepo,
		UserRepo:    userRepo,
		AsynqClient: asynqClient,
		Canceller:   canceller,
		Logger:      logger,
		RootDir:     rootDir,
	}
}

func (s *analysisService) getAnalysisFolderPath(
	userID, sampleID, analysisID uuid.UUID) string {
	return filepath.Join(
		s.RootDir,
		"uploads",
		"users", userID.String(),
		"samples", sampleID.String(),
		"analyses", analysisID.String(),
	)
}

func (s *analysisService) FindAll(ctx context.Context, userID uuid.UUID,
	filter models.AnalysisFilter, language string) (
		[]models.AnalysisResponse, error) {
	analyses, err := s.Repo.GetAnalyses(ctx, userID, filter)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "FindAll",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	responses := make([]models.AnalysisResponse, len(analyses))
	for i, analysis := range analyses {
		responses[i] = analysis.ToResponse(language)
	}

	return responses, nil
}

func (s *analysisService) FindManyByIDs(ctx context.Context,
	analysisIDs []uuid.UUID, userID uuid.UUID, language string) (
	[]models.AnalysisResponse, error) {
	if len(analysisIDs) > models.AnalysesByBatch {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "FindManyByIDs",
			logging.ExceededDownloadLimitError, ErrExceededDownloadLimit,
		)...)
		return nil, ErrExceededDownloadLimit
	}

	if len(analysisIDs) == 0 {
		return []models.AnalysisResponse{}, nil
	}

	analyses, err := s.Repo.GetAnalysesByIDs(ctx, analysisIDs, userID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "FindManyByIDs",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	var responses []models.AnalysisResponse
	for _, a := range analyses {
		responses = append(responses, a.ToResponse(language))
	}
	return responses, nil
}

func (s *analysisService) FindByID(ctx context.Context, analysisID,
	userID uuid.UUID, language string) (
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

	response := analysis.ToResponse(language)
	return &response, nil
}

func (s *analysisService) Create(ctx context.Context,
	input models.AnalysisCreateDTO, language string) (
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

	if sample.Fastq1 == nil && sample.Fastq2 == nil &&
		sample.Fasta == nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AnalysisService", "Create",
				logging.MissingFileError, ErrMissingFiles,
			)...)
		return nil, ErrMissingFiles
	}

	switch input.Type {
	case models.AnalysisTypeFastQC, models.AnalysisTypeComplete:
		if sample.Fastq1 == nil {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AnalysisService", "Create",
					logging.MissingFileError, ErrMissingFastq1,
				)...)
			return nil, ErrMissingFastq1
		}
		if sample.Fastq2 == nil {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AnalysisService", "Create",
					logging.MissingFileError, ErrMissingFastq2,
				)...)
			return nil, ErrMissingFastq2
		}
	case models.AnalysisTypeGenome:
		if (sample.Fastq1 == nil || sample.Fastq2 == nil) &&
			sample.Fasta == nil {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AnalysisService", "Create",
					logging.MissingFileError, ErrMissingFiles,
				)...)
			return nil, ErrMissingFiles
		}
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

	task, err := tasks.NewAnalysisProcessTask(analysis.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Create", logging.AsynqTaskError,
			err,
		)...)
	} else {
		info, err := s.AsynqClient.EnqueueContext(ctx, task,
			asynq.Queue(tasks.QueueAnalysis))
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AnalysisService", "Create",
				logging.RedisDispatchError, err,
			)...)
		} else {
			s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(
				"AnalysisService", "Create",
				logging.TaskEnqueuedSuccess, zap.String("task_id", info.ID),
				zap.String("queue", info.Queue),
			)...)

			taskID := info.ID
			analysis.TaskID = &taskID
			if err := s.Repo.UpdateAnalysis(ctx, &analysis); err != nil {
				s.Logger.Warn("Service Warning", logging.ServiceLogging(
					"AnalysisService", "Create",
					logging.DatabaseError, err,
				)...)
			}
		}
	}

	response := analysis.ToResponse(language)
	return &response, nil
}

func (s *analysisService) Update(ctx context.Context, analysisID uuid.UUID,
	input models.AdminAnalysisUpdateInput, language string) (
	*models.AnalysisResponse, error) {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Update", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "Update", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	if input.Status != nil &&
		!analysis.Status.CanTransitionTo(*input.Status) {
		return nil, ErrInvalidStatusTransition
	}

	validations.ApplyAnalysisUpdate(analysis, &input)

	if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AnalysisService", "Update",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	if analysis.Status == models.AnalysisStatusFailed &&
		analysis.TaskID != nil {
		if err := s.Canceller.CancelProcessing(*analysis.TaskID); err != nil {
			s.Logger.Warn("Service Warning", logging.ServiceLogging(
				"AnalysisService", "Update",
				logging.RedisDispatchError, err,
			)...)
		}
	}

	if analysis.Status == models.AnalysisStatusDone ||
		analysis.Status == models.AnalysisStatusFailed {
		task, err := tasks.NewAnalysisDoneEmailTask(analysis.ID)
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AnalysisService", "Update", logging.AsynqTaskError,
				err,
			)...)
		} else {
			info, err := s.AsynqClient.EnqueueContext(ctx, task,
				asynq.Queue(tasks.QueueEmail))
			if err != nil {
				s.Logger.Error("Service Error", logging.ServiceLogging(
					"AnalysisService", "Update",
					logging.RedisDispatchError, err,
				)...)
			} else {
				s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(
					"AnalysisService", "Update",
					logging.TaskEnqueuedSuccess, zap.String("task_id", info.ID),
					zap.String("queue", info.Queue),
				)...)
			}
		}
	}

	if analysis.Status == models.AnalysisStatusPending {
		task, err := tasks.NewAnalysisProcessTask(analysis.ID)
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AnalysisService", "Update", logging.AsynqTaskError,
				err,
			)...)
		} else {
			info, err := s.AsynqClient.EnqueueContext(ctx, task,
				asynq.Queue(tasks.QueueAnalysis))
			if err != nil {
				s.Logger.Error("Service Error", logging.ServiceLogging(
					"AnalysisService", "Update",
					logging.RedisDispatchError, err,
				)...)
			} else {
				s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(
					"AnalysisService", "Update",
					logging.TaskEnqueuedSuccess, zap.String("task_id", info.ID),
					zap.String("queue", info.Queue),
				)...)

				taskID := info.ID
				analysis.TaskID = &taskID
				if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
					s.Logger.Warn("Service Warning", logging.ServiceLogging(
						"AnalysisService", "Update",
						logging.DatabaseError, err,
					)...)
				}
			}
		}
	}

	response := analysis.ToResponse(language)
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

	uploadDir := s.getAnalysisFolderPath(analysis.UserID, analysis.SampleID,
		analysisID)
	if err := os.RemoveAll(uploadDir); err != nil {
		s.Logger.Warn("Service Warning", logging.ServiceLogging(
			"AnalysisService", "Delete", logging.DeleteFolderError, err,
		)...)
	}

	return nil
}

func (s *analysisService) DownloadZip(ctx context.Context, analysisID,
	userID uuid.UUID) (string, error) {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "DownloadZip", logging.DatabaseNotFoundError,
			err,
		)...)
		return "", ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "DownloadZip", logging.DatabaseError, err,
		)...)
		return "", ErrInternal
	}

	if userID != uuid.Nil && userID != analysis.UserID {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "DownloadZip", logging.Unauthorized, err,
		)...)
		return "", ErrUnauthorized
	}

	zipPath := ""
	if analysis.Status == models.AnalysisStatusDone &&
		analysis.ResultsZipPath != nil {
		zipPath = *analysis.ResultsZipPath
	}

	if zipPath == "" {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "DownloadZip", logging.MissingFileError,
			ErrZipNotFound,
		)...)
		return "", ErrZipNotFound
	}

	if _, err := os.Stat(zipPath); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "DownloadZip", logging.MissingFileError,
			ErrZipNotFound,
		)...)
		return "", ErrZipNotFound
	}

	return zipPath, nil
}

func (s *analysisService) DownloadBatchTSV(ctx context.Context,
	analysisIDs []uuid.UUID, userID uuid.UUID, language string) (
	[]models.AnalysisResponse, error) {
	if len(analysisIDs) > models.AnalysesByBatch {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "DownloadBatchTSV",
			logging.ExceededDownloadLimitError, ErrExceededDownloadLimit,
		)...)
		return nil, ErrExceededDownloadLimit
	}

	if len(analysisIDs) == 0 {
		return []models.AnalysisResponse{}, nil
	}

	analyses, err := s.Repo.GetAnalysesByIDs(ctx, analysisIDs, userID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisService", "DownloadBatchTSV",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	for _, a := range analyses {
		if a.Type == models.AnalysisTypeFastQC {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AnalysisService", "DownloadBatchTSV",
				logging.ExceededDownloadLimitError, ErrFastQCDownload,
			)...)
			return nil, ErrFastQCDownload
		}
	}

	var responses []models.AnalysisResponse
	for _, a := range analyses {
		responses = append(responses, a.ToResponse(language))
	}
	return responses, nil
}
