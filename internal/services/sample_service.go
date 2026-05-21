package services

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SampleService interface {
	PrepareSampleFolder(userID, sampleID uuid.UUID) (string, error)
	FindAll(ctx context.Context, input string, userID uuid.UUID,
		language string) ([]models.SampleResponse, error)
	FindByID(ctx context.Context, sampleID, userID uuid.UUID,
		language string) (*models.SampleResponse, error)
	Create(ctx context.Context, input models.SampleCreateInput,
		language string) (*models.SampleResponse, error)
	AttachFiles(ctx context.Context, sampleID, userID uuid.UUID,
		input models.SampleAttachmentInput) error
	Update(ctx context.Context, sampleID, userID uuid.UUID,
		input models.SampleUpdateInput,
		language string) (*models.SampleResponse, error)
	Delete(ctx context.Context, sampleID, userID uuid.UUID) error
}

type sampleService struct {
	Repo              repositories.SampleRepository
	CountryRepo       repositories.CountryRepository
	UserRepo          repositories.UserRepository
	OriginRepo        repositories.OriginRepository
	SampleSourceRepo  repositories.SampleSourceRepository
	MicroorganismRepo repositories.MicroorganismRepository
	SequencerRepo     repositories.SequencerRepository
	LaboratoryRepo    repositories.LaboratoryRepository
	HealthServiceRepo repositories.HealthServiceRepository
	RootDir           string
	Logger            *zap.Logger
}

func NewSampleService(
	repo repositories.SampleRepository,
	countryRepo repositories.CountryRepository,
	userRepo repositories.UserRepository,
	originRepo repositories.OriginRepository,
	sampleSourceRepo repositories.SampleSourceRepository,
	microorganismRepo repositories.MicroorganismRepository,
	sequencerRepo repositories.SequencerRepository,
	laboratoryRepo repositories.LaboratoryRepository,
	healthServiceRepo repositories.HealthServiceRepository,
	rootDir string,
	logger *zap.Logger) SampleService {
	return &sampleService{
		Repo:              repo,
		CountryRepo:       countryRepo,
		UserRepo:          userRepo,
		OriginRepo:        originRepo,
		SampleSourceRepo:  sampleSourceRepo,
		MicroorganismRepo: microorganismRepo,
		SequencerRepo:     sequencerRepo,
		LaboratoryRepo:    laboratoryRepo,
		HealthServiceRepo: healthServiceRepo,
		RootDir:           rootDir,
		Logger:            logger,
	}
}

func (s *sampleService) getSampleFolderPath(userID, sampleID uuid.UUID) string {
	return filepath.Join(
		s.RootDir,
		"uploads",
		"users", userID.String(),
		"samples", sampleID.String(),
	)
}

func (s *sampleService) PrepareSampleFolder(
	userID, sampleID uuid.UUID) (string, error) {
	basePath := s.getSampleFolderPath(userID, sampleID)

	if err := os.MkdirAll(basePath, 0755); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "PrepareSampleFolder",
			logging.CreateFolderError, err,
		)...)
		return "", ErrCreateFolder
	}

	return basePath, nil
}

func (s *sampleService) FindAll(ctx context.Context, input string,
	userID uuid.UUID, language string) ([]models.SampleResponse, error) {
	samples, err := s.Repo.GetSamples(ctx, input, userID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "FindAll",
			logging.DatabaseError, err,
		)...)
		return nil, ErrInternal
	}

	responses := make([]models.SampleResponse, len(samples))
	for i, sample := range samples {
		responses[i] = sample.ToResponse(language)
	}

	return responses, nil
}

func (s *sampleService) FindByID(
	ctx context.Context, sampleID, userID uuid.UUID,
	language string) (*models.SampleResponse, error) {
	sample, err := s.Repo.GetSampleByID(ctx, sampleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "FindByID", logging.DatabaseNotFoundError, err,
		)...)
		return nil, ErrNotFound
	}

	if userID != uuid.Nil && userID != sample.UserID {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "FindByID", logging.Unauthorized, err,
		)...)
		return nil, ErrUnauthorized
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "FindByID",
			logging.DatabaseError, err)...)
		return nil, ErrInternal
	}

	response := sample.ToResponse(language)
	return &response, nil
}

func (s *sampleService) Create(
	ctx context.Context,
	input models.SampleCreateInput,
	language string) (*models.SampleResponse, error) {
	country, err := s.CountryRepo.GetCountryByCode(ctx, input.CountryCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrInvalidCountryCode
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	user, err := s.UserRepo.GetUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrUserNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	origin, err := s.OriginRepo.GetOriginByID(ctx, input.OriginID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrOriginNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	sampleSource, err := s.SampleSourceRepo.GetSampleSourceByID(ctx,
		input.SampleSourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrSampleSourceNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	microorganism, err := s.MicroorganismRepo.GetMicroorganismByID(ctx,
		input.MicroorganismID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrMicroorganismNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	sequencer, err := s.SequencerRepo.GetSequencerByID(ctx,
		input.SequencerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrSequencerNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	laboratory, err := s.LaboratoryRepo.GetLaboratoryByID(ctx,
		input.LaboratoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrLaboratoryNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	healthService, err := s.HealthServiceRepo.GetHealthServiceByID(ctx,
		input.HealthServiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Create",
					logging.ExternalRepositoryNotFoundError, err,
				)...)
			return nil, ErrHealthServiceNotFound
		}
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	sample := models.Sample{
		Name:            input.Name,
		CollectionDate:  input.CollectionDate,
		RunNumber:       input.RunNumber,
		RunDate:         input.RunDate,
		City:            input.City,
		OriginCode:      input.OriginCode,
		Gender:          input.Gender,
		DateOfBirth:     input.DateOfBirth,
		CountryID:       country.ID,
		UserID:          user.ID,
		OriginID:        origin.ID,
		SampleSourceID:  sampleSource.ID,
		MicroorganismID: microorganism.ID,
		SequencerID:     sequencer.ID,
		LaboratoryID:    laboratory.ID,
		HealthServiceID: healthService.ID,
	}

	sample.Country = *country
	sample.User = *user
	sample.Origin = *origin
	sample.SampleSource = *sampleSource
	sample.Microorganism = *microorganism
	sample.Sequencer = *sequencer
	sample.Laboratory = *laboratory
	sample.HealthService = *healthService

	if err := s.Repo.CreateSample(ctx, &sample); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Create",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := sample.ToResponse(language)
	return &response, nil
}

func (s *sampleService) AttachFiles(ctx context.Context,
	sampleID, userID uuid.UUID,
	input models.SampleAttachmentInput) error {
	sample, err := s.Repo.GetSampleByID(ctx, sampleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "AttachFiles", logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "AttachFiles", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if userID != uuid.Nil && userID != sample.UserID {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "AttachFiles", logging.Unauthorized, err,
		)...)
		return ErrUnauthorized
	}

	if input.Fastq1 == nil && input.Fastq2 == nil && input.Fasta == nil {
		return ErrMissingFiles
	} else if input.Fastq1 != nil && input.Fastq2 == nil {
		return ErrMissingFastq2
	} else if input.Fastq1 == nil && input.Fastq2 != nil {
		return ErrMissingFastq1
	}

	oldFastq1 := sample.Fastq1
	oldFastq2 := sample.Fastq2
	oldFasta := sample.Fasta

	validations.ApplySampleFilesUpdate(sample, &input)

	if err := s.Repo.UpdateSample(ctx, sample); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "AttachFiles",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if oldFastq1 != nil && input.Fastq1 != nil && *oldFastq1 != *input.Fastq1 {
		if err := os.Remove(*oldFastq1); err != nil {
			s.Logger.Warn("Service Warning", logging.ServiceLogging(
				"SampleService", "AttachFiles", logging.DeleteFileError, err,
			)...)
		}
	}
	if oldFastq2 != nil && input.Fastq2 != nil && *oldFastq2 != *input.Fastq2 {
		if err := os.Remove(*oldFastq2); err != nil {
			s.Logger.Warn("Service Warning", logging.ServiceLogging(
				"SampleService", "AttachFiles", logging.DeleteFileError, err,
			)...)
		}
	}
	if oldFasta != nil && input.Fasta != nil && *oldFasta != *input.Fasta {
		if err := os.Remove(*oldFasta); err != nil {
			s.Logger.Warn("Service Warning", logging.ServiceLogging(
				"SampleService", "AttachFiles", logging.DeleteFileError, err,
			)...)
		}
	}

	return nil
}

func (s *sampleService) Update(
	ctx context.Context, sampleID, userID uuid.UUID,
	input models.SampleUpdateInput,
	language string) (*models.SampleResponse, error) {
	existingSample, err := s.Repo.GetSampleByID(ctx, sampleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Update",
				logging.DatabaseNotFoundError, err,
			)...)
		return nil, ErrNotFound
	}
	if err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Update",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	if userID != uuid.Nil && userID != existingSample.UserID {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "Update", logging.Unauthorized, err,
		)...)
		return nil, ErrUnauthorized
	}

	if input.CountryCode != nil {
		country, err := s.CountryRepo.GetCountryByCode(ctx, *input.CountryCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrInvalidCountryCode
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.CountryID = country.ID
		existingSample.Country = *country
	}

	if input.UserID != nil {
		user, err := s.UserRepo.GetUserByID(ctx, *input.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrUserNotFound
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.UserID = user.ID
		existingSample.User = *user
	}

	if input.OriginID != nil {
		origin, err := s.OriginRepo.GetOriginByID(ctx, *input.OriginID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrOriginNotFound
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.OriginID = origin.ID
		existingSample.Origin = *origin
	}

	if input.SampleSourceID != nil {
		sampleSource, err := s.SampleSourceRepo.GetSampleSourceByID(ctx,
			*input.SampleSourceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrSampleSourceNotFound
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.SampleSourceID = sampleSource.ID
		existingSample.SampleSource = *sampleSource
	}

	if input.MicroorganismID != nil {
		microorganism, err := s.MicroorganismRepo.GetMicroorganismByID(ctx,
			*input.MicroorganismID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrMicroorganismNotFound
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.MicroorganismID = microorganism.ID
		existingSample.Microorganism = *microorganism
	}

	if input.SequencerID != nil {
		sequencer, err := s.SequencerRepo.GetSequencerByID(ctx,
			*input.SequencerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrSequencerNotFound
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.SequencerID = sequencer.ID
		existingSample.Sequencer = *sequencer
	}

	if input.LaboratoryID != nil {
		laboratory, err := s.LaboratoryRepo.GetLaboratoryByID(ctx,
			*input.LaboratoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrLaboratoryNotFound
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.LaboratoryID = laboratory.ID
		existingSample.Laboratory = *laboratory
	}

	if input.HealthServiceID != nil {
		healthService, err := s.HealthServiceRepo.GetHealthServiceByID(ctx,
			*input.HealthServiceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.Error("Service Error",
					logging.ServiceLogging(
						"SampleService", "Update",
						logging.ExternalRepositoryNotFoundError, err,
					)...)
				return nil, ErrHealthServiceNotFound
			}
			s.Logger.Error("Service Error",
				logging.ServiceLogging(
					"SampleService", "Update",
					logging.ExternalRepositoryError, err,
				)...)
			return nil, ErrInternal
		}
		existingSample.HealthServiceID = healthService.ID
		existingSample.HealthService = *healthService
	}

	validations.ApplySampleUpdate(existingSample, &input)

	if err := s.Repo.UpdateSample(ctx, existingSample); err != nil {
		s.Logger.Error("Service Error",
			logging.ServiceLogging(
				"SampleService", "Update",
				logging.DatabaseError, err,
			)...)
		return nil, ErrInternal
	}

	response := existingSample.ToResponse(language)
	return &response, nil
}

func (s *sampleService) Delete(ctx context.Context,
	sampleID, userID uuid.UUID) error {
	sample, err := s.Repo.GetSampleByID(ctx, sampleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "Delete", logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}

	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	if userID != uuid.Nil && userID != sample.UserID {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "Delete", logging.Unauthorized, err,
		)...)
		return ErrUnauthorized
	}

	if err := s.Repo.DeleteSample(ctx, sample); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"SampleService", "Delete", logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	uploadDir := s.getSampleFolderPath(sample.UserID, sampleID)
	if err := os.RemoveAll(uploadDir); err != nil {
		s.Logger.Warn("Service Warning", logging.ServiceLogging(
			"SampleService", "Delete", logging.DeleteFolderError, err,
		)...)
	}

	return nil
}
