package services_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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

func TestAnalysisFindAll(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return []models.Analysis{mock}, nil
			},
		}

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil, t.TempDir())
		result, err := svc.FindAll(ctx, uuid.Nil)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mock.ToResponse(), result[0])
	})

	t.Run("Error", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysesFunc: func(ctx context.Context,
				userID uuid.UUID) ([]models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		result, err := svc.FindAll(ctx, uuid.Nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAnalysisFindManyByIDs(t *testing.T) {
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil, t.TempDir())
		result, err := svc.FindManyByIDs(ctx, []uuid.UUID{mock.ID},
			mock.User.ID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mock.ToResponse(), result[0])
	})

	t.Run("Success - Empty Analysis IDs", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil, t.TempDir())
		result, err := svc.FindManyByIDs(ctx, []uuid.UUID{},
			mock.User.ID)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Error - Exceeded Limit", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}

		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		result, err := svc.FindManyByIDs(ctx, make([]uuid.UUID,
			models.AnalysesByBatch+1), mock.User.ID)

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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		result, err := svc.FindManyByIDs(ctx, []uuid.UUID{mock.ID},
			mock.User.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAnalysisFindByID(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	t.Run("Success", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
		}

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil, t.TempDir())
		result, err := svc.FindByID(ctx, mock.ID, mock.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		expected := mock.ToResponse()
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		result, err := svc.FindByID(ctx, mock.ID, mock.UserID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		result, err := svc.FindByID(ctx, mock.ID, uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return nil, services.ErrInternal
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		result, err := svc.FindByID(ctx, mock.ID, mock.UserID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAnalysisCreate(t *testing.T) {
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

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			userRepo, enqueuer, mockLogger, t.TempDir())
		result, err := svc.Create(ctx, input)

		expected := models.AnalysisResponse{
			Type:     input.Type,
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

	t.Run("Success - Soft Fail Asynq", func(t *testing.T) {
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

		failingEnqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
				return nil, errors.New("redis timeout")
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			userRepo, failingEnqueuer, mockLogger, t.TempDir())
		result, err := svc.Create(ctx, input)

		expected := models.AnalysisResponse{
			Type:     input.Type,
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

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			nil, nil, mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			nil, nil, mockLogger, t.TempDir())
		result, err := svc.Create(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sample Missing Fastq1 in FASTQC analysis",
		func(t *testing.T) {
			analysisRepo := &mocks.MockAnalysisRepository{}
			sampleRepo := &mocks.MockSampleRepository{
				GetSampleByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Sample, error) {
					sample := &models.Sample{}
					return sample, nil
				},
			}

			mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

			svc := services.NewAnalysisService(analysisRepo, sampleRepo,
				nil, nil, mockLogger, t.TempDir())

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

	t.Run("Error - Sample Missing Fastq2 in FASTQC analysis",
		func(t *testing.T) {
			analysisRepo := &mocks.MockAnalysisRepository{}
			sampleRepo := &mocks.MockSampleRepository{
				GetSampleByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Sample, error) {
					fastq1 := "reads1.fastq"
					sample := &models.Sample{Fastq1: &fastq1}
					return sample, nil
				},
			}

			mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

			svc := services.NewAnalysisService(analysisRepo, sampleRepo,
				nil, nil, mockLogger, t.TempDir())

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

	t.Run("Success - Change Analysis Type from Complete To Genome",
		func(t *testing.T) {
			analysisRepo := &mocks.MockAnalysisRepository{}
			sampleRepo := &mocks.MockSampleRepository{
				GetSampleByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Sample, error) {
					fasta := "assembly.fasta"
					sample := &models.Sample{
						ID:    mock.Sample.ID,
						Name:  mock.Sample.Name,
						Fasta: &fasta,
					}
					return sample, nil
				},
			}
			userRepo := &mocks.MockUserRepository{
				GetUserByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.User, error) {
					return &mock.User, nil
				},
			}

			mockLogger, logs := testutils.NewMockLogger(zap.InfoLevel)

			enqueuer := &mocks.MockTaskEnqueuer{}
			svc := services.NewAnalysisService(analysisRepo, sampleRepo,
				userRepo, enqueuer, mockLogger, t.TempDir())

			result, err := svc.Create(ctx, input)

			expected := models.AnalysisResponse{
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

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			userRepo, nil, mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			userRepo, nil, mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, sampleRepo, userRepo,
			nil, mockLogger, t.TempDir())
		result, err := svc.Create(ctx, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAnalysisUpdate(t *testing.T) {
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil,
			enqueuer, mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil,
			enqueuer, mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil,
			failingEnqueuer, mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger, t.TempDir())
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger, t.TempDir())
		result, err := svc.Update(ctx, mock.ID, updateInputRunning)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAnalysisDelete(t *testing.T) {
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil, t.TempDir())
		err := svc.Delete(ctx, mock.ID, mock.UserID)

		assert.NoError(t, err)
	})

	t.Run("Success - Deletes Analysis Folder", func(t *testing.T) {
		rootDir := t.TempDir()
		analysisFolder := filepath.Join(rootDir, "uploads", "users",
			mock.UserID.String(), "samples", mock.SampleID.String(),
			"analyses", mock.ID.String())
		err := os.MkdirAll(analysisFolder, 0755)
		assert.NoError(t, err)
		_, err = os.Create(filepath.Join(analysisFolder, "result.zip"))
		assert.NoError(t, err)

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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil,
			rootDir)
		err = svc.Delete(ctx, mock.ID, mock.UserID)

		assert.NoError(t, err)
		_, statErr := os.Stat(analysisFolder)
		assert.True(t, os.IsNotExist(statErr))
	})

	t.Run("Success - Folder Deletion Failure Logs Warning", func(t *testing.T) {
		rootDirFile := filepath.Join(t.TempDir(), "root.txt")
		err := os.WriteFile(rootDirFile, []byte("x"), 0644)
		assert.NoError(t, err)

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

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)
		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger, rootDirFile)
		err = svc.Delete(ctx, mock.ID, mock.UserID)

		assert.NoError(t, err)
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		err := svc.Delete(ctx, mock.ID, mock.UserID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return &mock, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		err := svc.Delete(ctx, mock.ID, uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		err := svc.Delete(ctx, mock.ID, mock.UserID)

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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		err := svc.Delete(ctx, runningMock.ID, runningMock.UserID)

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

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		err := svc.Delete(ctx, mock.ID, mock.UserID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})
}
