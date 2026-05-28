package services_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

func TestPrepareSampleFolder(t *testing.T) {
	mock := testmodels.CreateMockSample()
	rootDir := t.TempDir()

	expected := filepath.Join(rootDir, "uploads", "users",
		mock.UserID.String(), "samples", mock.ID.String())

	sampleRepo := &mocks.MockSampleRepository{}
	svc := services.NewSampleService(sampleRepo, nil, nil, nil,
		nil, nil, nil, nil, nil, rootDir, nil)
	result, err := svc.PrepareSampleFolder(mock.UserID, mock.ID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSampleFindAll(t *testing.T) {
	mock := testmodels.CreateMockSample()

	t.Run("Success", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID) ([]models.Sample, error) {
				return []models.Sample{mock}, nil
			},
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, nil, nil, nil, nil, t.TempDir(), nil)
		result, err := svc.FindAll(context.Background(), "", uuid.Nil, "en")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mock.ToResponse("en"), result[0])
	})

	t.Run("Error", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSamplesFunc: func(ctx context.Context, input string,
				userID uuid.UUID) ([]models.Sample, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil,
			nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.FindAll(context.Background(), "", uuid.Nil, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Empty(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleFindByID(t *testing.T) {
	mock := testmodels.CreateMockSample()

	t.Run("Success", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), nil)
		result, err := svc.FindByID(context.Background(), mock.ID, uuid.Nil,
			"en")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		expected := mock.ToResponse("en")
		assert.Equal(t, expected, *result)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil,
			nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.FindByID(context.Background(), uuid.New(),
			uuid.Nil, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.FindByID(context.Background(), mock.ID,
			uuid.New(), "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.FindByID(context.Background(), uuid.New(),
			uuid.Nil, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleCreate(t *testing.T) {
	mock := testmodels.CreateMockSample()
	input := testmodels.NewSampleCreateDTO(mock)

	happyRepos := func() (
		*mocks.MockCountryRepository,
		*mocks.MockUserRepository,
		*mocks.MockOriginRepository,
		*mocks.MockSampleSourceRepository,
		*mocks.MockMicroorganismRepository,
		*mocks.MockSequencerRepository,
		*mocks.MockLaboratoryRepository,
		*mocks.MockHealthServiceRepository,
	) {
		return &mocks.MockCountryRepository{
				GetCountryByCodeFunc: func(ctx context.Context,
					code string) (*models.Country, error) {
					c := mock.Country
					return &c, nil
				},
			},
			&mocks.MockUserRepository{
				GetUserByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.User, error) {
					u := mock.User
					return &u, nil
				},
			},
			&mocks.MockOriginRepository{
				GetOriginByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Origin, error) {
					o := mock.Origin
					return &o, nil
				},
			},
			&mocks.MockSampleSourceRepository{
				GetSampleSourceByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.SampleSource, error) {
					ss := mock.SampleSource
					return &ss, nil
				},
			},
			&mocks.MockMicroorganismRepository{
				GetMicroorganismByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Microorganism, error) {
					m := mock.Microorganism
					return &m, nil
				},
			},
			&mocks.MockSequencerRepository{
				GetSequencerByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Sequencer, error) {
					s := mock.Sequencer
					return &s, nil
				},
			},
			&mocks.MockLaboratoryRepository{
				GetLaboratoryByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.Laboratory, error) {
					l := mock.Laboratory
					return &l, nil
				},
			},
			&mocks.MockHealthServiceRepository{
				GetHealthServiceByIDFunc: func(ctx context.Context,
					ID uuid.UUID) (*models.HealthService, error) {
					hs := mock.HealthService
					return &hs, nil
				},
			}
	}

	t.Run("Success", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			CreateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}
		countryRepo, userRepo, originRepo, ssRepo, microRepo, seqRepo,
			labRepo, hsRepo := happyRepos()

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, labRepo, hsRepo,
			t.TempDir(), nil)
		result, err := svc.Create(context.Background(), input, "en")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mock.Name, result.Name)
	})

	t.Run("Error - Invalid Country Code", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				code string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, nil, nil,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCountryCode)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				code string) (*models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, nil, nil,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, _, _, _, _, _, _, _ := happyRepos()
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			nil, nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUserNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - User Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, _, _, _, _, _, _, _ := happyRepos()
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			nil, nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Origin Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, _, _, _, _, _, _ := happyRepos()
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrOriginNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Origin Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, _, _, _, _, _, _ := happyRepos()
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - SampleSource Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, _, _, _, _, _ := happyRepos()
		ssRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrSampleSourceNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - SampleSource Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, _, _, _, _, _ := happyRepos()
		ssRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Microorganism Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, _, _, _, _ := happyRepos()
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, nil, nil, nil, t.TempDir(),
			mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMicroorganismNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Microorganism Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, _, _, _, _ := happyRepos()
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, nil, nil, nil, t.TempDir(),
			mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sequencer Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, microRepo,
			_, _, _ := happyRepos()
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, nil, nil, t.TempDir(),
			mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrSequencerNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sequencer Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, microRepo, _, _,
			_ := happyRepos()
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, nil, nil, t.TempDir(),
			mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Laboratory Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, microRepo, seqRepo,
			_, _ := happyRepos()
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, labRepo, nil,
			t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrLaboratoryNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Laboratory Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, microRepo,
			seqRepo, _, _ := happyRepos()
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, labRepo, nil, t.TempDir(),
			mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - HealthService Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, microRepo, seqRepo,
			labRepo, _ := happyRepos()
		hsRepo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, labRepo, hsRepo,
			t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrHealthServiceNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - HealthService Database", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{}
		countryRepo, userRepo, originRepo, ssRepo, microRepo,
			seqRepo, labRepo, _ := happyRepos()
		hsRepo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, labRepo, hsRepo,
			t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			CreateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return gorm.ErrInvalidTransaction
			},
		}
		countryRepo, userRepo, originRepo, ssRepo, microRepo, seqRepo, labRepo,
			hsRepo := happyRepos()

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, userRepo,
			originRepo, ssRepo, microRepo, seqRepo, labRepo, hsRepo,
			t.TempDir(), mockLogger)
		result, err := svc.Create(context.Background(), input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleAttachFiles(t *testing.T) {
	mock := testmodels.CreateMockSample()

	fastq1 := "new_read1.fastq"
	fastq2 := "new_read2.fastq"
	fasta := "assembly.fasta"

	t.Run("Success - Fastq pair", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		mockLogger, _ := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fastq1: &fastq1, Fastq2: &fastq2})

		assert.NoError(t, err)
	})

	t.Run("Success - Fasta only", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				s := mock
				s.Fastq1 = nil
				s.Fastq2 = nil
				return &s, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		mockLogger, _ := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fasta: &fasta})

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), uuid.New(), uuid.Nil,
			models.SampleAttachmentInput{Fastq1: &fastq1, Fastq2: &fastq2})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.New(),
			models.SampleAttachmentInput{Fastq1: &fastq1, Fastq2: &fastq2})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Missing All Files", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), nil)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFiles)
	})

	t.Run("Error - Missing Fastq2", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), nil)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fastq1: &fastq1})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFastq2)
	})

	t.Run("Error - Missing Fastq1", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), nil)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fastq2: &fastq2})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMissingFastq1)
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fastq1: &fastq1, Fastq2: &fastq2})

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Fastq1 not deleted", func(t *testing.T) {
		oldFastq1 := "/nonexistent/old_read1.fastq"
		sampleWithOldFiles := mock
		sampleWithOldFiles.Fastq1 = &oldFastq1

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &sampleWithOldFiles, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fastq1: &fastq1, Fastq2: &fastq2})

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.FilterLevelExact(zapcore.WarnLevel).Len())
	})

	t.Run("Error - Fastq2 not deleted", func(t *testing.T) {
		oldFastq2 := "/nonexistent/old_read2.fastq"
		sampleWithOldFiles := mock
		sampleWithOldFiles.Fastq2 = &oldFastq2

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &sampleWithOldFiles, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fastq1: &fastq1, Fastq2: &fastq2})

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.FilterLevelExact(zapcore.WarnLevel).Len())
	})

	t.Run("Error - Fasta not deleted", func(t *testing.T) {
		oldFasta := "/nonexistent/old_assembly.fasta"
		sampleWithOldFiles := mock
		sampleWithOldFiles.Fasta = &oldFasta

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &sampleWithOldFiles, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.AttachFiles(context.Background(), mock.ID, uuid.Nil,
			models.SampleAttachmentInput{Fasta: &fasta})

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.FilterLevelExact(zapcore.WarnLevel).Len())
	})
}

func TestSampleUpdate(t *testing.T) {
	mock := testmodels.CreateMockSample()

	newName := "Updated Sample"
	input := models.SampleUpdateDTO{Name: &newName}

	t.Run("Success", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), nil)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			input, "en")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newName, result.Name)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), uuid.New(), uuid.Nil,
			input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Get by ID DB Internal", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.New(),
			input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Not Found", func(t *testing.T) {
		code := "XX"
		inputWithCountry := models.SampleUpdateDTO{CountryCode: &code}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				c string) (*models.Country, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, nil, nil,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithCountry, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidCountryCode)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country Database", func(t *testing.T) {
		code := "XX"
		inputWithCountry := models.SampleUpdateDTO{CountryCode: &code}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		countryRepo := &mocks.MockCountryRepository{
			GetCountryByCodeFunc: func(ctx context.Context,
				c string) (*models.Country, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, countryRepo, nil, nil,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithCountry, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		userID := uuid.New()
		inputWithUser := models.SampleUpdateDTO{UserID: &userID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, userRepo, nil,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithUser, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUserNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - User Database", func(t *testing.T) {
		userID := uuid.New()
		inputWithUser := models.SampleUpdateDTO{UserID: &userID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		userRepo := &mocks.MockUserRepository{
			GetUserByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, userRepo, nil,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithUser, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Origin Not Found", func(t *testing.T) {
		originID := uuid.New()
		inputWithOrigin := models.SampleUpdateDTO{OriginID: &originID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, originRepo,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithOrigin, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrOriginNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Origin Database", func(t *testing.T) {
		originID := uuid.New()
		inputWithOrigin := models.SampleUpdateDTO{OriginID: &originID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		originRepo := &mocks.MockOriginRepository{
			GetOriginByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, originRepo,
			nil, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithOrigin, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - SampleSource Not Found", func(t *testing.T) {
		ssID := uuid.New()
		inputWithSS := models.SampleUpdateDTO{SampleSourceID: &ssID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		ssRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			ssRepo, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithSS, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrSampleSourceNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - SampleSource Database", func(t *testing.T) {
		ssID := uuid.New()
		inputWithSS := models.SampleUpdateDTO{SampleSourceID: &ssID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		ssRepo := &mocks.MockSampleSourceRepository{
			GetSampleSourceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.SampleSource, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			ssRepo, nil, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithSS, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Microorganism Not Found", func(t *testing.T) {
		microID := uuid.New()
		inputWithMicro := models.SampleUpdateDTO{MicroorganismID: &microID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, microRepo, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithMicro, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrMicroorganismNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Microorganism Database", func(t *testing.T) {
		microID := uuid.New()
		inputWithMicro := models.SampleUpdateDTO{MicroorganismID: &microID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		microRepo := &mocks.MockMicroorganismRepository{
			GetMicroorganismByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Microorganism, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, microRepo, nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithMicro, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sequencer Not Found", func(t *testing.T) {
		seqID := uuid.New()
		inputWithSeq := models.SampleUpdateDTO{SequencerID: &seqID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, nil, seqRepo, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithSeq, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrSequencerNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Sequencer Database", func(t *testing.T) {
		seqID := uuid.New()
		inputWithSeq := models.SampleUpdateDTO{SequencerID: &seqID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		seqRepo := &mocks.MockSequencerRepository{
			GetSequencerByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sequencer, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, nil, seqRepo, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithSeq, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Laboratory Not Found", func(t *testing.T) {
		labID := uuid.New()
		inputWithLab := models.SampleUpdateDTO{LaboratoryID: &labID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, nil, nil, labRepo, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithLab, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrLaboratoryNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Laboratory Database", func(t *testing.T) {
		labID := uuid.New()
		inputWithLab := models.SampleUpdateDTO{LaboratoryID: &labID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		labRepo := &mocks.MockLaboratoryRepository{
			GetLaboratoryByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, nil, nil, labRepo, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithLab, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - HealthService Not Found", func(t *testing.T) {
		hsID := uuid.New()
		inputWithHS := models.SampleUpdateDTO{HealthServiceID: &hsID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		hsRepo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, nil, nil, nil, hsRepo, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithHS, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrHealthServiceNotFound)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - HealthService Database", func(t *testing.T) {
		hsID := uuid.New()
		inputWithHS := models.SampleUpdateDTO{HealthServiceID: &hsID}

		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}
		hsRepo := &mocks.MockHealthServiceRepository{
			GetHealthServiceByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.HealthService, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil,
			nil, nil, nil, nil, hsRepo, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			inputWithHS, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Update DB Internal", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
			UpdateSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		result, err := svc.Update(context.Background(), mock.ID, uuid.Nil,
			input, "en")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, result)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestSampleDelete(t *testing.T) {
	mock := testmodels.CreateMockSample()

	t.Run("Success", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
			DeleteSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), nil)
		err := svc.Delete(context.Background(), mock.ID, uuid.Nil)

		assert.NoError(t, err)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.Delete(context.Background(), uuid.New(), uuid.Nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.Delete(context.Background(), mock.ID, uuid.New())

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrUnauthorized)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal", func(t *testing.T) {
		sampleRepo := mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
			DeleteSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return gorm.ErrInvalidTransaction
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewSampleService(&sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, t.TempDir(), mockLogger)
		err := svc.Delete(context.Background(), mock.ID, uuid.Nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Delete Folder", func(t *testing.T) {
		sampleRepo := &mocks.MockSampleRepository{
			GetSampleByIDFunc: func(ctx context.Context,
				ID uuid.UUID) (*models.Sample, error) {
				return &mock, nil
			},
			DeleteSampleFunc: func(ctx context.Context,
				sample *models.Sample) error {
				return nil
			},
		}

		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)
		var mockRootDir string
		switch runtime.GOOS {
		case "windows":
			mockRootDir = "\\Windows\\System32"
		case "darwin":
			mockRootDir = "/var/root"
		default:
			mockRootDir = "/root"
		}

		svc := services.NewSampleService(sampleRepo, nil, nil, nil, nil, nil,
			nil, nil, nil, mockRootDir, mockLogger)
		err := svc.Delete(context.Background(), mock.ID, uuid.Nil)

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.FilterLevelExact(zapcore.WarnLevel).Len())
	})
}
