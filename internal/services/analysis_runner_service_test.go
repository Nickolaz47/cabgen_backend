package services_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestAnalysisRunnerRun(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{}
		enqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(_ context.Context,
				task *asynq.Task, _ ...asynq.Option) (*asynq.TaskInfo,
				error) {
				assert.Equal(t, tasks.TaskTypeAnalysisDoneEmail,
					task.Type())
				return &asynq.TaskInfo{ID: "t1", Queue: "emails"}, nil
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, enqueuer,
			zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusDone, updated.Status)
		assert.NotNil(t, updated.StartedAt)
		assert.NotNil(t, updated.FinishedAt)
		assert.Nil(t, updated.ErrorMessage)
		assert.NotNil(t, updated.FastQC1)
		assert.NotNil(t, updated.FastQC2)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisRunnerService(repo, nil,
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, uuid.New())

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal on Get", func(t *testing.T) {
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisRunnerService(repo, nil,
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, uuid.New())

		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB Internal on Update", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				_ *models.Analysis) error {
				return gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisRunnerService(repo, nil,
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - FastQC", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunFastQCFunc: func(_ context.Context, read1, read2,
				outputDir string) (string, string, error) {
				return "", "", errors.New("fastqc crashed")
			},
		}
		mockLogger, _ := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrFastQC.Error())
		assert.NotContains(t, *updated.ErrorMessage, "fastqc crashed")
	})

	t.Run("Error - Unknown Type", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = "NONSENSE"
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}

		svc := services.NewAnalysisRunnerService(repo,
			&mocks.MockCabgenPipeline{}, &mocks.MockTaskEnqueuer{},
			zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrUnknownAnalysisType.Error())
	})

	t.Run("Error - Prepare Folders", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}

		svc := services.NewAnalysisRunnerService(repo,
			&mocks.MockCabgenPipeline{}, &mocks.MockTaskEnqueuer{},
			zap.NewNop(), "/nonexistent_root_no_perms/x")
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrPrepareFolders.Error())
	})

	t.Run("Error - Unicycler", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = nil

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunUnicyclerFunc: func(_ context.Context, threads int,
				read1, read2, spadesPath, outputDir string) (
				string, error) {
				return "", errors.New("spades missing")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrUnicycler.Error())
		assert.NotContains(t, *updated.ErrorMessage,
			"spades missing")
	})

	t.Run("Error - Enqueue", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				_ *models.Analysis) error {
				return nil
			},
		}
		enqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(_ context.Context,
				_ *asynq.Task, _ ...asynq.Option) (*asynq.TaskInfo,
				error) {
				return nil, errors.New("redis down")
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisRunnerService(repo,
			&mocks.MockCabgenPipeline{}, enqueuer, mockLogger,
			t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestAnalysisRunnerGenome(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Existing FASTA", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		fasta := "contigs.fasta"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fasta

		unicyclerCalled := false
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				_ *models.Analysis) error {
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunUnicyclerFunc: func(_ context.Context, threads int,
				read1, read2, spadesPath, outputDir string) (
				string, error) {
				unicyclerCalled = true
				return "assembly.fa", nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return os.WriteFile(outputFile,
					[]byte("placeholder\n"), 0644)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.False(t, unicyclerCalled,
			"Unicycler should not run when Fasta already present")
	})

	t.Run("Success - No FASTA", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = nil

		unicyclerCalled := false
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				_ *models.Analysis) error {
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunUnicyclerFunc: func(_ context.Context, threads int,
				read1, read2, spadesPath, outputDir string) (
				string, error) {
				unicyclerCalled = true
				return "assembly.fa", nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return os.WriteFile(outputFile,
					[]byte("placeholder\n"), 0644)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.True(t, unicyclerCalled)
	})

	t.Run("Error - Abricate", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		fasta := "contigs.fasta"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fasta

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}
		var failedDB string
		pl := &mocks.MockCabgenPipeline{
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				if failedDB == "" {
					failedDB = db
				}
				return errors.New("abricate segfault")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrAbricate.Error())
		assert.NotEmpty(t, failedDB,
			"at least one DB should have failed")
		assert.NotContains(t, *updated.ErrorMessage,
			"abricate segfault")
	})

	t.Run("Error - CheckM", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		fasta := "contigs.fasta"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fasta

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunCheckMFunc: func(_ context.Context, threads int,
				sample, assemblyDir, outputDir string) (
				*pipeline.CheckMResult, error) {
				return nil, fmt.Errorf("checkm db corrupt")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrCheckM.Error())
		assert.NotContains(t, *updated.ErrorMessage,
			"checkm db corrupt")
	})
}

func TestAnalysisRunnerComplete(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeComplete
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		fasta := "contigs.fasta"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fasta

		fastqcCalled := false
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				_ *models.Analysis) error {
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunFastQCFunc: func(_ context.Context, read1, read2,
				outputDir string) (string, string, error) {
				fastqcCalled = true
				return "qc1.html", "qc2.html", nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return os.WriteFile(outputFile,
					[]byte("placeholder\n"), 0644)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.True(t, fastqcCalled,
			"Complete should call FastQC first")
	})

	t.Run("Error - FastQC", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeComplete
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		fasta := "contigs.fasta"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fasta

		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunFastQCFunc: func(_ context.Context, read1, read2,
				outputDir string) (string, string, error) {
				return "", "", errors.New("fastqc timeout")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl,
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrFastQC.Error())
	})
}

func TestAnalysisRunnerPrepareFolders(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		root := t.TempDir()
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				_ *models.Analysis) error {
				return nil
			},
		}

		svc := services.NewAnalysisRunnerService(repo,
			&mocks.MockCabgenPipeline{}, &mocks.MockTaskEnqueuer{},
			zap.NewNop(), root)
		err := svc.Run(context.Background(), mock.ID)

		assert.NoError(t, err)
		for _, sub := range []string{"qc", "assembly", "amr", "report"} {
			assert.DirExists(t,
				filepath.Join(root, mock.ID.String(), sub))
		}
	})
}
