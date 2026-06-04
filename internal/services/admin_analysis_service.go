package services

import (
	"context"
	"errors"

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

type AdminAnalysisService interface {
	FindAll(ctx context.Context) ([]models.AnalysisAdminResponse, error)
	FindManyByIDs(ctx context.Context, analysisIDs []uuid.UUID) (
		[]models.AnalysisAdminResponse, error)
	FindByID(ctx context.Context, analysisID uuid.UUID) (
		*models.AnalysisAdminResponse, error)
	Create(ctx context.Context, input models.AnalysisCreateDTO) (
		*models.AnalysisAdminResponse, error)
	Update(ctx context.Context, analysisID uuid.UUID,
		input models.AdminAnalysisUpdateInput) (
		*models.AnalysisAdminResponse, error)
	Delete(ctx context.Context, analysisID uuid.UUID) error
}

type adminAnalysisService struct {
	Repo        repositories.AnalysisRepository
	SampleRepo  repositories.SampleRepository
	UserRepo    repositories.UserRepository
	AsynqClient TaskEnqueuer
	Logger      *zap.Logger
}

func NewAdminAnalysisService(
	repo repositories.AnalysisRepository,
	sampleRepo repositories.SampleRepository,
	userRepo repositories.UserRepository,
	AsynqClient TaskEnqueuer,
	logger *zap.Logger,
) AdminAnalysisService {
	return &adminAnalysisService{
		Repo:        repo,
		SampleRepo:  sampleRepo,
		UserRepo:    userRepo,
		AsynqClient: AsynqClient,
		Logger:      logger,
	}
}

func (s *adminAnalysisService) FindAll(ctx context.Context) (
	[]models.AnalysisAdminResponse, error) {
	analyses, err := s.Repo.GetAnalyses(ctx, uuid.Nil)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "FindAll",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	responses := make([]models.AnalysisAdminResponse, len(analyses))
	for i, analysis := range analyses {
		responses[i] = analysis.ToAdminResponse()
	}

	return responses, nil
}

func (s *adminAnalysisService) FindManyByIDs(ctx context.Context,
	analysisIDs []uuid.UUID) (
	[]models.AnalysisAdminResponse, error) {
	if len(analysisIDs) > models.AnalysesByBatch {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "FindManyByIDs",
			logging.ExceededDownloadLimitError, ErrExceededDownloadLimit,
		)...)
		return nil, ErrExceededDownloadLimit
	}

	if len(analysisIDs) == 0 {
		return []models.AnalysisAdminResponse{}, nil
	}

	analyses, err := s.Repo.GetAnalysesByIDs(ctx, analysisIDs, uuid.Nil)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "FindManyByIDs",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	var responses []models.AnalysisAdminResponse
	for _, a := range analyses {
		responses = append(responses, a.ToAdminResponse())
	}
	return responses, nil
}

func (s *adminAnalysisService) FindByID(ctx context.Context,
	analysisID uuid.UUID) (*models.AnalysisAdminResponse, error) {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "FindByID", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "FindByID",
			logging.DatabaseError, err)...)
		return nil, ErrInternal
	}

	response := analysis.ToAdminResponse()
	return &response, nil
}

func (s *adminAnalysisService) Create(ctx context.Context,
	input models.AnalysisCreateDTO) (
	*models.AnalysisAdminResponse, error) {
	sample, err := s.SampleRepo.GetSampleByID(ctx, input.SampleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminAnalysisService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrSampleNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminAnalysisService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	if input.Type == models.AnalysisTypeFastQC &&
		sample.Fastq1 == nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminAnalysisService", "Create",
				logging.MissingFileError, ErrMissingFastq1,
			)...)
		return nil, ErrMissingFastq1
	} else if input.Type == models.AnalysisTypeFastQC &&
		sample.Fastq2 == nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminAnalysisService", "Create",
				logging.MissingFileError, ErrMissingFastq2,
			)...)
		return nil, ErrMissingFastq2
	}

	if input.Type == models.AnalysisTypeComplete && (sample.Fastq1 == nil || sample.Fastq2 == nil) &&
		sample.Fasta != nil {
		input.Type = models.AnalysisTypeGenome
	}

	user, err := s.UserRepo.GetUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"AdminAnalysisService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrUserNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminAnalysisService", "Create",
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
				"AdminAnalysisService", "Create",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	task, err := tasks.NewProcessAnalysisTask(analysis.ID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "Create", logging.AsynqTaskError,
			err,
		)...)
	} else {
		info, err := s.AsynqClient.EnqueueContext(ctx, task,
			asynq.Queue(tasks.QueueAnalysis))
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AdminAnalysisService", "Create",
				logging.RedisDispatchError, err,
			)...)
		} else {
			s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(
				"AdminAnalysisService", "Create",
				logging.TaskEnqueuedSuccess, zap.String("task_id", info.ID),
				zap.String("queue", info.Queue),
			)...)
		}
	}

	response := analysis.ToAdminResponse()
	return &response, nil
}

func (s *adminAnalysisService) Update(ctx context.Context, analysisID uuid.UUID,
	input models.AdminAnalysisUpdateInput) (
	*models.AnalysisAdminResponse, error) {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "Update", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "Update", logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	validations.ApplyAnalysisUpdate(analysis, &input)

	if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"AdminAnalysisService", "Update",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	if analysis.Status == models.AnalysisStatusDone ||
		analysis.Status == models.AnalysisStatusFailed {
		task, err := tasks.NewAnalysisDoneEmailTask(analysis.ID)
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AdminAnalysisService", "Update", logging.AsynqTaskError,
				err,
			)...)
		} else {
			info, err := s.AsynqClient.EnqueueContext(ctx, task,
				asynq.Queue(tasks.QueueEmail))
			if err != nil {
				s.Logger.Error("Service Error", logging.ServiceLogging(
					"AdminAnalysisService", "Update",
					logging.RedisDispatchError, err,
				)...)
			} else {
				s.Logger.Info("Redis Task Info", logging.ServiceInfoLogging(
					"AdminAnalysisService", "Update",
					logging.TaskEnqueuedSuccess, zap.String("task_id", info.ID),
					zap.String("queue", info.Queue),
				)...)
			}
		}
	}

	response := analysis.ToAdminResponse()
	return &response, nil
}

func (s *adminAnalysisService) Delete(ctx context.Context,
	analysisID uuid.UUID) error {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "Delete", logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if analysis.Status == models.AnalysisStatusRunning {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrDeleteRunningAnalysis
	}

	if err := s.Repo.DeleteAnalysis(ctx, analysis); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AdminAnalysisService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	return nil
}
