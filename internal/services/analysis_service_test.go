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
		result, err := svc.FindAll(ctx, uuid.Nil, "en")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mock.ToResponse("en"), result[0])
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
		result, err := svc.FindAll(ctx, uuid.Nil, "en")

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
			mock.User.ID, "en")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mock.ToResponse("en"), result[0])
	})

	t.Run("Success - Empty Analysis IDs", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, nil, t.TempDir())
		result, err := svc.FindManyByIDs(ctx, []uuid.UUID{},
			mock.User.ID, "en")

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Error - Exceeded Limit", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}

		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil, mockLogger, t.TempDir())
		result, err := svc.FindManyByIDs(ctx, make([]uuid.UUID,
			models.AnalysesByBatch+1), mock.User.ID, "en")

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
			mock.User.ID, "en")

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
		result, err := svc.FindByID(ctx, mock.ID, mock.UserID, "en")

		assert.NoError(t, err)
		assert.NotNil(t, result)

		expected := mock.ToResponse("en")
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
		result, err := svc.FindByID(ctx, mock.ID, mock.UserID, "en")

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
		result, err := svc.FindByID(ctx, mock.ID, uuid.New(), "en")

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
		result, err := svc.FindByID(ctx, mock.ID, mock.UserID, "en")

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
		result, err := svc.Create(ctx, input, "en")

		expected := models.AnalysisResponse{
			Type:     input.Type,
			Status:   models.AnalysisStatusPending,
			Sample:   mock.Sample.OriginCode,
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
		result, err := svc.Create(ctx, input, "en")

		expected := models.AnalysisResponse{
			Type:     input.Type,
			Status:   models.AnalysisStatusPending,
			Sample:   mock.Sample.OriginCode,
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
		result, err := svc.Create(ctx, input, "en")

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
		result, err := svc.Create(ctx, input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - No Files At All", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &models.Sample{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			nil, nil, mockLogger, t.TempDir())

		result, err := svc.Create(ctx, models.AnalysisCreateDTO{
			Type:     models.AnalysisTypeComplete,
			SampleID: input.SampleID,
			UserID:   input.UserID,
		}, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFiles)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sample Missing Fastq1 in FASTQC analysis",
		func(t *testing.T) {
			analysisRepo := &mocks.MockAnalysisRepository{}
			sampleRepo := &mocks.MockSampleRepository{
				GetSampleByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Sample, error) {
					fasta := "assembly.fasta"
					sample := &models.Sample{Fasta: &fasta}
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
			result, err := svc.Create(ctx, errorInput, "en")

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
					fasta := "assembly.fasta"
					sample := &models.Sample{Fastq1: &fastq1, Fasta: &fasta}
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
			result, err := svc.Create(ctx, errorInput, "en")

			assert.Error(t, err)
			assert.ErrorIs(t, err, services.ErrMissingFastq2)
			assert.Nil(t, result)
			assert.Equal(t, 1, logs.Len())
		})

	t.Run("Error - Complete With Only Fasta",
		func(t *testing.T) {
			analysisRepo := &mocks.MockAnalysisRepository{}
			sampleRepo := &mocks.MockSampleRepository{
				GetSampleByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Sample, error) {
					fasta := "assembly.fasta"
					sample := &models.Sample{
						ID:         mock.Sample.ID,
						OriginCode: mock.Sample.OriginCode,
						Fasta:      &fasta,
					}
					return sample, nil
				},
			}

			mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

			svc := services.NewAnalysisService(analysisRepo, sampleRepo,
				nil, nil, mockLogger, t.TempDir())

			result, err := svc.Create(ctx, input, "en")

			assert.Error(t, err)
			assert.ErrorIs(t, err, services.ErrMissingFastq1)
			assert.Nil(t, result)
			assert.Equal(t, 1, logs.Len())
		})

	t.Run("Error - Complete Missing Fastq1", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &models.Sample{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			nil, nil, mockLogger, t.TempDir())

		result, err := svc.Create(ctx, models.AnalysisCreateDTO{
			Type:     models.AnalysisTypeComplete,
			SampleID: input.SampleID,
			UserID:   input.UserID,
		}, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFiles)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Complete Missing Fastq2", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				fastq1 := "reads1.fastq"
				return &models.Sample{Fastq1: &fastq1}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			nil, nil, mockLogger, t.TempDir())

		result, err := svc.Create(ctx, models.AnalysisCreateDTO{
			Type:     models.AnalysisTypeComplete,
			SampleID: input.SampleID,
			UserID:   input.UserID,
		}, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFastq2)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Genome With Only Fasta", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				fasta := "assembly.fasta"
				return &models.Sample{
					ID:         mock.Sample.ID,
					OriginCode: mock.Sample.OriginCode,
					Fasta:      &fasta,
				}, nil
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

		result, err := svc.Create(ctx, models.AnalysisCreateDTO{
			Type:     models.AnalysisTypeGenome,
			SampleID: input.SampleID,
			UserID:   input.UserID,
		}, "en")

		expected := models.AnalysisResponse{
			Type:     models.AnalysisTypeGenome,
			Status:   models.AnalysisStatusPending,
			Sample:   mock.Sample.OriginCode,
			SampleID: mock.Sample.ID,
			User:     mock.User.Username,
			UserID:   mock.User.ID,
		}

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, *result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Genome With No Files", func(t *testing.T) {
		analysisRepo := &mocks.MockAnalysisRepository{}
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &models.Sample{}, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			nil, nil, mockLogger, t.TempDir())

		result, err := svc.Create(ctx, models.AnalysisCreateDTO{
			Type:     models.AnalysisTypeGenome,
			SampleID: input.SampleID,
			UserID:   input.UserID,
		}, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFiles)
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

		svc := services.NewAnalysisService(analysisRepo, sampleRepo,
			userRepo, nil, mockLogger, t.TempDir())
		result, err := svc.Create(ctx, input, "en")

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
		result, err := svc.Create(ctx, input, "en")

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
		result, err := svc.Create(ctx, input, "en")

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
		result, err := svc.Update(ctx, mock.ID, updateInputRunning, "en")

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
		result, err := svc.Update(ctx, mock.ID, updateInputDone, "en")

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
		result, err := svc.Update(ctx, mock.ID, updateInputFailed, "en")

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
		result, err := svc.Update(ctx, mock.ID, updateInputDone, "en")

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
		result, err := svc.Update(ctx, mock.ID, updateInputRunning, "en")

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
		result, err := svc.Update(ctx, mock.ID, updateInputRunning, "en")

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
		result, err := svc.Update(ctx, mock.ID, updateInputRunning, "en")

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

func TestAnalysisDownloadZip(t *testing.T) {
	ctx := context.Background()

	newRepo := func(analysisReturn func() (*models.Analysis, error)) *mocks.MockAnalysisRepository {
		return &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(ctx context.Context,
				analysisID uuid.UUID) (*models.Analysis, error) {
				return analysisReturn()
			},
		}
	}

	t.Run("Success", func(t *testing.T) {
		rootDir := t.TempDir()
		zipPath := filepath.Join(rootDir, "results.zip")
		err := os.WriteFile(zipPath, []byte("x"), 0644)
		assert.NoError(t, err)

		mock := testmodels.CreateMockAnalysis()
		mock.Status = models.AnalysisStatusDone
		mock.ResultsZipPath = &zipPath

		svc := services.NewAnalysisService(newRepo(func() (*models.Analysis,
			error) {
			return &mock, nil
		}), nil, nil, nil, zap.NewNop(),
			rootDir)
		gotPath, err := svc.DownloadZip(ctx, mock.ID, mock.UserID)

		assert.NoError(t, err)
		assert.Equal(t, zipPath, gotPath)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(newRepo(func() (*models.Analysis,
			error) {
			return nil, gorm.ErrRecordNotFound
		}), nil, nil, nil,
			mockLogger, t.TempDir())
		_, err := svc.DownloadZip(ctx, uuid.New(), uuid.Nil)

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal on Get", func(t *testing.T) {
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(newRepo(func() (*models.Analysis,
			error) {
			return nil, gorm.ErrInvalidTransaction
		}), nil, nil, nil,
			mockLogger, t.TempDir())
		_, err := svc.DownloadZip(ctx, uuid.New(), uuid.Nil)

		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.ResultsZipPath = nil
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(newRepo(func() (*models.Analysis,
			error) {
			return &mock, nil
		}), nil, nil, nil, mockLogger,
			t.TempDir())
		_, err := svc.DownloadZip(ctx, mock.ID, uuid.New())

		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Analysis Not Done", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Status = models.AnalysisStatusRunning
		mock.ResultsZipPath = nil
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(newRepo(func() (*models.Analysis,
			error) {
			return &mock, nil
		}), nil, nil, nil, mockLogger,
			t.TempDir())
		_, err := svc.DownloadZip(ctx, mock.ID, mock.UserID)

		assert.ErrorIs(t, err, services.ErrZipNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - ResultsZipPath Not Set", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Status = models.AnalysisStatusDone
		mock.ResultsZipPath = nil
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(newRepo(func() (*models.Analysis,
			error) {
			return &mock, nil
		}), nil, nil, nil, mockLogger,
			t.TempDir())
		_, err := svc.DownloadZip(ctx, mock.ID, mock.UserID)

		assert.ErrorIs(t, err, services.ErrZipNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Zip File Missing on Disk", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Status = models.AnalysisStatusDone
		missingPath := filepath.Join(t.TempDir(), "missing.zip")
		mock.ResultsZipPath = &missingPath
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisService(newRepo(func() (*models.Analysis,
			error) {
			return &mock, nil
		}), nil, nil, nil, mockLogger,
			t.TempDir())
		_, err := svc.DownloadZip(ctx, mock.ID, mock.UserID)

		assert.ErrorIs(t, err, services.ErrZipNotFound)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAnalysisDownloadBatchTSV(t *testing.T) {
	ctx := context.Background()
	mock := testmodels.CreateMockAnalysis()

	fastqcMock := testmodels.CreateMockAnalysis()
	fastqcMock.Type = models.AnalysisTypeFastQC

	analysisRepo := &mocks.MockAnalysisRepository{
		GetAnalysesByIDsFunc: func(ctx context.Context, ids []uuid.UUID,
			userID uuid.UUID) ([]models.Analysis, error) {
			return []models.Analysis{mock, fastqcMock}, nil
		},
	}

	t.Run("Success", func(t *testing.T) {
		successRepo := &mocks.MockAnalysisRepository{
			GetAnalysesByIDsFunc: func(ctx context.Context, ids []uuid.UUID,
				userID uuid.UUID) ([]models.Analysis, error) {
				return []models.Analysis{mock}, nil
			},
		}
		svc := services.NewAnalysisService(successRepo, nil, nil, nil,
			zap.NewNop(), t.TempDir())
		responses, err := svc.DownloadBatchTSV(ctx,
			[]uuid.UUID{mock.ID}, mock.UserID, "en")

		assert.NoError(t, err)
		assert.Len(t, responses, 1)
	})

	t.Run("Error - Exceeded Limit", func(t *testing.T) {
		ids := make([]uuid.UUID, models.AnalysesByBatch+1)
		for i := range ids {
			ids[i] = uuid.New()
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger, t.TempDir())
		responses, err := svc.DownloadBatchTSV(ctx, ids, mock.UserID, "en")

		assert.ErrorIs(t, err, services.ErrExceededDownloadLimit)
		assert.Nil(t, responses)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Empty IDs Returns Empty List", func(t *testing.T) {
		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil,
			zap.NewNop(), t.TempDir())
		responses, err := svc.DownloadBatchTSV(ctx, []uuid.UUID{},
			mock.UserID, "en")

		assert.NoError(t, err)
		assert.Empty(t, responses)
	})

	t.Run("Error - FASTQC in Batch", func(t *testing.T) {
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAnalysisService(analysisRepo, nil, nil, nil,
			mockLogger, t.TempDir())
		responses, err := svc.DownloadBatchTSV(ctx,
			[]uuid.UUID{mock.ID, fastqcMock.ID}, mock.UserID, "en")

		assert.ErrorIs(t, err, services.ErrFastQCDownload)
		assert.Nil(t, responses)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		failRepo := &mocks.MockAnalysisRepository{
			GetAnalysesByIDsFunc: func(ctx context.Context, ids []uuid.UUID,
				userID uuid.UUID) ([]models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)
		svc := services.NewAnalysisService(failRepo, nil, nil, nil,
			mockLogger, t.TempDir())
		responses, err := svc.DownloadBatchTSV(ctx,
			[]uuid.UUID{mock.ID}, mock.UserID, "en")

		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, responses)
		assert.Equal(t, 1, logs.Len())
	})
}
