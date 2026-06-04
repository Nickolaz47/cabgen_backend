package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

func TestAdminAnalysisFindAll(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return []models.Analysis{mock}, nil
			},
		}

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, nil)
		result, err := svc.FindAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mock.ToAdminResponse(), result[0])
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		result, err := svc.FindAll(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminAnalysisFindManyByIDs(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesByIDsFunc: func(ctx context.Context,
				analysisIDs []uuid.UUID, userID uuid.UUID) (
				[]models.Analysis, error) {
				return []models.Analysis{mock}, nil
			},
		}

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, nil)
		result, err := svc.FindManyByIDs(ctx, []uuid.UUID{mock.ID})

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mock.ToAdminResponse(), result[0])
	})

	t.Run("Success - Empty Analysis IDs", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, nil)
		result, err := svc.FindManyByIDs(ctx, []uuid.UUID{})

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Error - Exceeded Limit", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}

		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, mockLogger)
		result, err := svc.FindManyByIDs(ctx, make([]uuid.UUID,
			models.AnalysesByBatch+1))

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrExceededDownloadLimit)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesByIDsFunc: func(ctx context.Context,
				analysisIDs []uuid.UUID, userID uuid.UUID) (
				[]models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, mockLogger)
		result, err := svc.FindManyByIDs(ctx, []uuid.UUID{mock.ID})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminAnalysisFindByID(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
		}

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, nil)
		result, err := svc.FindByID(ctx, mock.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		expected := mock.ToAdminResponse()
		assert.Equal(t, expected, *result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		result, err := svc.FindByID(ctx, mock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		result, err := svc.FindByID(ctx, mock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminAnalysisCreate(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()
	input := testmodels.NewAnalysisCreateDTO(mock)

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock.Sample, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &mock.User, nil
			},
		}

		enqueuer := &mocks.MockTaskEnqueuer{}
		mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)
		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo,
			userRepo, enqueuer, mockLogger)
		result, err := svc.Create(ctx, input)

		expected := models.AnalysisAdminResponse{
			Type:     input.Type,
			Status:   models.AnalysisStatusPending,
			Sample:   mock.Sample.Name,
			SampleID: mock.Sample.ID,
			User:     mock.Sample.User.Username,
			UserID:   mock.Sample.User.ID,
		}

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, *result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Change Analysis Type from Complete To Genome",
		func(t *testing.T) {
			analysisRepo := &mocks.MockAnalysisRepository{}
			sampleRepo := &mocks.MockSampleRepository{
				GetSampleByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Sample, error) {
					fasta := "assembly.fasta"
					return &models.Sample{
						ID:    mock.Sample.ID,
						Name:  mock.Sample.Name,
						Fasta: &fasta,
					}, nil
				},
			}
			userRepo := &mocks.MockUserRepository{
				GetUserByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.User, error) {
					return &mock.User, nil
				},
			}

			enqueuer := &mocks.MockTaskEnqueuer{}
			mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)

			svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo,
				userRepo, enqueuer, mockLogger)
			result, err := svc.Create(ctx, input)

			expected := models.AnalysisAdminResponse{
				Type:     models.AnalysisTypeGenome,
				Status:   models.AnalysisStatusPending,
				Sample:   mock.Sample.Name,
				SampleID: mock.Sample.ID,
				User:     mock.User.Username,
				UserID:   mock.User.ID,
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expected, *result)
			assert.Equal(t, 1, logs.Len())
		})

	t.Run("Error - Sample Not Found", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo, nil,
			nil, mockLogger)
		result, err := svc.Create(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrSampleNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sample Database", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo, nil,
			nil, mockLogger)
		result, err := svc.Create(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Missing Fastq1 in FastQC Analysis", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &models.Sample{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo, nil,
			nil, mockLogger)

		errorInput := models.AnalysisCreateDTO{
			Type:     models.AnalysisTypeFastQC,
			SampleID: input.SampleID,
			UserID:   input.UserID,
		}
		result, err := svc.Create(ctx, errorInput)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFastq1)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Missing Fastq2 in FastQC Analysis", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				fastq1 := "reads1.fastq"
				return &models.Sample{Fastq1: &fastq1}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo, nil,
			nil, mockLogger)

		errorInput := models.AnalysisCreateDTO{
			Type:     models.AnalysisTypeFastQC,
			SampleID: input.SampleID,
			UserID:   input.UserID,
		}
		result, err := svc.Create(ctx, errorInput)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFastq2)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock.Sample, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo,
			userRepo, nil, mockLogger)
		result, err := svc.Create(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUserNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - User Database", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock.Sample, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo,
			userRepo, nil, mockLogger)
		result, err := svc.Create(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			CreateAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return gorm.ErrInvalidTransaction
			},
		}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock.Sample, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return &mock.User, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, sampleRepo,
			userRepo, nil, mockLogger)
		result, err := svc.Create(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminAnalysisUpdate(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	statusRunning := models.AnalysisStatusRunning
	updateInputRunning := models.AdminAnalysisUpdateInput{
		Status: &statusRunning,
	}

	statusDone := models.AnalysisStatusDone
	updateInputDone := models.AdminAnalysisUpdateInput{
		Status: &statusDone,
	}

	statusFailed := models.AnalysisStatusFailed
	updateInputFailed := models.AdminAnalysisUpdateInput{
		Status: &statusFailed,
	}

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
			UpdateAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return nil
			},
		}

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, nil)
		result, err := svc.Update(ctx, mock.ID, updateInputRunning)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.AnalysisStatusRunning, result.Status)
		assert.NotNil(t, result.StartedAt)
	})

	t.Run("Success - Status Done Enqueues Email Task", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
			UpdateAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return nil
			},
		}

		enqueuer := &mocks.MockTaskEnqueuer{}
		mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil,
			enqueuer, mockLogger)
		result, err := svc.Update(ctx, mock.ID, updateInputDone)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.AnalysisStatusDone, result.Status)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Status Failed Enqueues Email Task", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
			UpdateAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return nil
			},
		}

		enqueuer := &mocks.MockTaskEnqueuer{}
		mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil,
			enqueuer, mockLogger)
		result, err := svc.Update(ctx, mock.ID, updateInputFailed)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.AnalysisStatusFailed, result.Status)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Soft Fail Asynq Enqueue Error", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
			UpdateAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return nil
			},
		}

		failingEnqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task,
				opts ...asynq.Option) (*asynq.TaskInfo, error) {
				return nil, errors.New("redis timeout")
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil,
			failingEnqueuer, mockLogger)
		result, err := svc.Update(ctx, mock.ID, updateInputDone)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.AnalysisStatusDone, result.Status)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		result, err := svc.Update(ctx, mock.ID, updateInputRunning)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal on Get", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		result, err := svc.Update(ctx, mock.ID, updateInputRunning)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal on Update", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
			UpdateAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		result, err := svc.Update(ctx, mock.ID, updateInputRunning)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAdminAnalysisDelete(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
			DeleteAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return nil
			},
		}

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil, nil)
		err := svc.Delete(ctx, mock.ID)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		err := svc.Delete(ctx, mock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal on Get", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		err := svc.Delete(ctx, mock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Delete Running Analysis", func(t *testing.T) {
		runningMock := testmodels.CreateMockAnalysis()
		runningMock.Status = models.AnalysisStatusRunning

		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &runningMock, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		err := svc.Delete(ctx, runningMock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrDeleteRunningAnalysis)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal on Delete", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
			DeleteAnalysisFunc: func(ctx context.Context,
				analysis *models.Analysis) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAdminAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger)
		err := svc.Delete(ctx, mock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
