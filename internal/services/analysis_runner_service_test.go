package services_test

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
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

	originalConcurrency := config.AnalysisConcurrency
	config.AnalysisConcurrency = 4
	t.Cleanup(func() { config.AnalysisConcurrency = originalConcurrency })

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

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{}, enqueuer,
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
		assert.Empty(t, updated.Step)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisRunnerService(repo, nil, &mocks.MockCommander{},
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

		svc := services.NewAnalysisRunnerService(repo, nil, &mocks.MockCommander{},
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

		svc := services.NewAnalysisRunnerService(repo, nil, &mocks.MockCommander{},
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
		var steps []models.AnalysisStep
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				steps = append(steps, analysis.Step)
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

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrFastQC.Error())
		assert.NotContains(t, *updated.ErrorMessage, "fastqc crashed")
		assert.Contains(t, steps, models.StepFastQC)
		assert.Empty(t, updated.Step)
	})

	t.Run("Error - FastQC Input Error Preserved", func(t *testing.T) {
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
				return "", "", pipeline.ErrCorruptedInput
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrCorruptedInput.Error())
		assert.NotContains(t, *updated.ErrorMessage,
			pipeline.ErrFastQC.Error())
	})

	t.Run("Error - Unicycler Input Error Preserved", func(t *testing.T) {
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
				read1, read2, spadesPath, outputDir, outputFile string) (
				string, error) {
				return "", pipeline.ErrInvalidFormat
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrInvalidFormat.Error())
		assert.NotContains(t, *updated.ErrorMessage,
			pipeline.ErrUnicycler.Error())
	})

	t.Run("Error - Prokka", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fasta := filepath.Join(t.TempDir(), "contigs.fasta")
		os.WriteFile(fasta, []byte(">seq1\nATCG\n"), 0644)
		mock.Sample.Fastq1 = nil
		mock.Sample.Fastq2 = nil
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
			RunProkkaFunc: func(_ context.Context, threads int,
				assembly, outputDir string) error {
				return errors.New("prokka crashed")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrProkka.Error())
		assert.NotContains(t, *updated.ErrorMessage, "prokka crashed")
	})

	t.Run("Error - Kraken2", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fasta := filepath.Join(t.TempDir(), "contigs.fasta")
		os.WriteFile(fasta, []byte(">seq1\nATCG\n"), 0644)
		mock.Sample.Fastq1 = nil
		mock.Sample.Fastq2 = nil
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
			RunKraken2Func: func(_ context.Context, threads int,
				assembly, outputDir string) (*pipeline.KrakenSpecies,
				*pipeline.KrakenSpecies, error) {
				return nil, nil, errors.New("kraken2 crashed")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrKraken2.Error())
		assert.NotContains(t, *updated.ErrorMessage, "kraken2 crashed")
	})

	t.Run("Error - Species", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fasta := filepath.Join(t.TempDir(), "contigs.fasta")
		os.WriteFile(fasta, []byte(">seq1\nATCG\n"), 0644)
		mock.Sample.Fastq1 = nil
		mock.Sample.Fastq2 = nil
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
			RunKraken2Func: func(_ context.Context, threads int,
				assembly, outputDir string) (*pipeline.KrakenSpecies,
				*pipeline.KrakenSpecies, error) {
				return &pipeline.KrakenSpecies{
					Name: "Escherichia coli", Count: 100,
				}, nil, nil
			},
			ProcessSpeciesFunc: func(_ context.Context, threads int,
				sampleID, mostCommon, assemblyPath, outputDir string) (
				*pipeline.SpeciesResult, error) {
				return nil, errors.New("species crashed")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrSpecies.Error())
		assert.NotContains(t, *updated.ErrorMessage, "species crashed")
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

		emailCalls := 0
		enqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(_ context.Context, task *asynq.Task,
				_ ...asynq.Option) (*asynq.TaskInfo, error) {
				emailCalls++
				return nil, nil
			},
		}

		svc := services.NewAnalysisRunnerService(repo,
			&mocks.MockCabgenPipeline{}, &mocks.MockCommander{}, enqueuer,
			zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrUnknownAnalysisType.Error())
		assert.Equal(t, 0, emailCalls,
			"email should not be enqueued on non-final failure")
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
			&mocks.MockCabgenPipeline{}, &mocks.MockCommander{}, &mocks.MockTaskEnqueuer{},
			zap.NewNop(), "/nonexistent_root_no_perms/x")
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrPrepareFolders.Error())
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
		var steps []models.AnalysisStep
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				steps = append(steps, analysis.Step)
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunUnicyclerFunc: func(_ context.Context, threads int,
				read1, read2, spadesPath, outputDir, outputFile string) (
				string, error) {
				return "", errors.New("spades missing")
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrUnicycler.Error())
		assert.NotContains(t, *updated.ErrorMessage,
			"spades missing")
		assert.Contains(t, steps, models.StepUnicycler)
		assert.Empty(t, updated.Step)
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
			&mocks.MockCabgenPipeline{}, &mocks.MockCommander{}, enqueuer, mockLogger,
			t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.Equal(t, 1, logs.Len())
	})
}

func writeAbricateOutput(outputFile string) error {
	line := "file\tseq1\t1\t100\t+\tgene\t1-100/100\t=\t0/0\t95.0\t98.0\tresfinder\tAF123\tproduct"
	return os.WriteFile(outputFile, []byte(line+"\n"), 0644)
}

func newResfinderRef(t *testing.T) string {
	t.Helper()
	ref := strings.Repeat("P", 10000) + "\n" +
		strings.Join([]string{"blaTEM", "", "", "", "", "", "", "",
			"", "", "", "", "", "", "", "", "Ampicillin"}, "\t") + "\n"
	path := filepath.Join(t.TempDir(), "resfinder_ref.txt")
	if err := os.WriteFile(path, []byte(ref), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAnalysisRunnerGenome(t *testing.T) {
	ctx := context.Background()

	originalConcurrency := config.AnalysisConcurrency
	config.AnalysisConcurrency = 4
	t.Cleanup(func() { config.AnalysisConcurrency = originalConcurrency })

	t.Run("Success - Existing FASTA", func(t *testing.T) {
		tmpDir := t.TempDir()
		fastaPath := filepath.Join(tmpDir, "contigs.fasta")
		err := os.WriteFile(fastaPath, []byte(">seq1\nATCGATCG\n"), 0644)
		assert.NoError(t, err)

		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fastaPath

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
			Config: pipeline.ToolsConfig{
				ResfinderDBPath: newResfinderRef(t),
			},
			RunUnicyclerFunc: func(_ context.Context, threads int,
				read1, read2, spadesPath, outputDir, outputFile string) (
				string, error) {
				unicyclerCalled = true
				return "assembly.fa", nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return writeAbricateOutput(outputFile)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err = svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.False(t, unicyclerCalled,
			"Unicycler should not run when Fasta already present")
	})

	t.Run("Success - FASTA Only Copies to AssemblyDir", func(t *testing.T) {
		tmpDir := t.TempDir()
		fastaContent := ">seq1\nATCGATCG\n"
		fastaPath := filepath.Join(tmpDir, "user_genome.fasta")
		err := os.WriteFile(fastaPath, []byte(fastaContent), 0644)
		assert.NoError(t, err)

		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		mock.Sample.Fastq1 = nil
		mock.Sample.Fastq2 = nil
		mock.Sample.Fasta = &fastaPath

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
			Config: pipeline.ToolsConfig{
				ResfinderDBPath: newResfinderRef(t),
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return writeAbricateOutput(outputFile)
			},
		}

		rootDir := t.TempDir()
		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), rootDir)
		err = svc.Run(ctx, mock.ID)

		assert.NoError(t, err)

		assemblyDir := filepath.Join(rootDir, "uploads", "users",
			mock.UserID.String(), "samples", mock.SampleID.String(),
			"analyses", mock.ID.String(), "assembly")
		copiedFasta := filepath.Join(assemblyDir, "user_genome.fasta")
		assert.FileExists(t, copiedFasta)

		content, err := os.ReadFile(copiedFasta)
		assert.NoError(t, err)
		assert.Equal(t, fastaContent, string(content))
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
		var capturedOutputFile string
		var persistedSample *models.Sample
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
			UpdateSampleFunc: func(_ context.Context,
				sample *models.Sample) error {
				persistedSample = sample
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			Config: pipeline.ToolsConfig{
				ResfinderDBPath: newResfinderRef(t),
			},
			RunUnicyclerFunc: func(_ context.Context, threads int,
				read1, read2, spadesPath, outputDir, outputFile string) (
				string, error) {
				unicyclerCalled = true
				capturedOutputFile = outputFile
				return "assembly.fa", nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return writeAbricateOutput(outputFile)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.True(t, unicyclerCalled)
		assert.NotNil(t, persistedSample)
		assert.Equal(t, "assembly.fa", *persistedSample.Fasta)
		assert.Equal(t, mock.Sample.OriginCode+"_assembly.fasta",
			capturedOutputFile)
	})

	t.Run("Success - FASTA not found, fallback to reads", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		fasta := "/nonexistent/path/genome.fasta"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fasta

		unicyclerCalled := false
		var persistedSample *models.Sample
		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)
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
			UpdateSampleFunc: func(_ context.Context,
				sample *models.Sample) error {
				persistedSample = sample
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			Config: pipeline.ToolsConfig{
				ResfinderDBPath: newResfinderRef(t),
			},
			RunUnicyclerFunc: func(_ context.Context, threads int,
				read1, read2, spadesPath, outputDir, outputFile string) (
				string, error) {
				unicyclerCalled = true
				return "assembly.fa", nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return writeAbricateOutput(outputFile)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.True(t, unicyclerCalled,
			"Unicycler should run when FASTA file is missing")
		assert.NotNil(t, persistedSample)
		assert.Equal(t, "assembly.fa", *persistedSample.Fasta)

		found := false
		for _, entry := range logs.All() {
			if strings.Contains(entry.Message, "FASTA file not found") {
				found = true
				break
			}
		}
		assert.True(t, found, "should warn about missing FASTA file")
	})

	t.Run("Error - No input files", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		mock.Sample.Fastq1 = nil
		mock.Sample.Fastq2 = nil
		mock.Sample.Fasta = nil

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
			&mocks.MockCabgenPipeline{}, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
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
		var steps []models.AnalysisStep
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				steps = append(steps, analysis.Step)
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

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrAbricate.Error())
		assert.NotEmpty(t, failedDB,
			"at least one DB should have failed")
		assert.NotContains(t, *updated.ErrorMessage,
			"abricate segfault")
		assert.Contains(t, steps, models.StepAbricate)
		assert.Empty(t, updated.Step)
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
		var steps []models.AnalysisStep
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				steps = append(steps, analysis.Step)
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

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrCheckM.Error())
		assert.NotContains(t, *updated.ErrorMessage,
			"checkm db corrupt")
		assert.Contains(t, steps, models.StepCheckM)
		assert.Empty(t, updated.Step)
	})

	t.Run("Error - ProcessResfinder", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeGenome
		mock.Status = models.AnalysisStatusPending
		fq1, fq2 := "r1.fq", "r2.fq"
		fasta := "contigs.fasta"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2
		mock.Sample.Fasta = &fasta

		updated := (*models.Analysis)(nil)
		var steps []models.AnalysisStep
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				updated = analysis
				steps = append(steps, analysis.Step)
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				line := "file\tseq1\t1\t100\t+\tgene\t1-100/100\t=\t0/0\t95.0\t98.0\tresfinder\tAF123\tproduct"
				return os.WriteFile(outputFile,
					[]byte(line+"\n"), 0644)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, steps, models.StepAbricate)
		assert.Empty(t, updated.Step)
	})

	t.Run("Warning - CalculateCoverage", func(t *testing.T) {
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
			Config: pipeline.ToolsConfig{
				ResfinderDBPath: newResfinderRef(t),
			},
			RunCheckMFunc: func(_ context.Context, threads int,
				sample, assemblyDir, outputDir string) (
				*pipeline.CheckMResult, error) {
				return &pipeline.CheckMResult{
					Completeness: "99.5", Contamination: "0.5",
					GenomeSize: "5000000", N50: "100000",
				}, nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return writeAbricateOutput(outputFile)
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.ErrorLevel)

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusDone, updated.Status)
		assert.GreaterOrEqual(t, logs.Len(), 1)
	})

	t.Run("Success - Secondary Species Written When Contamination > 5%",
		func(t *testing.T) {
			mock := testmodels.CreateMockAnalysis()
			mock.Type = models.AnalysisTypeGenome
			mock.Status = models.AnalysisStatusPending
			fasta := filepath.Join(t.TempDir(), "contigs.fasta")
			os.WriteFile(fasta, []byte(">seq1\nATCG\n"), 0644)
			mock.Sample.Fastq1 = nil
			mock.Sample.Fastq2 = nil
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
				Config: pipeline.ToolsConfig{
					ResfinderDBPath: newResfinderRef(t),
				},
				RunCheckMFunc: func(_ context.Context, threads int,
					sample, assemblyDir, outputDir string) (
					*pipeline.CheckMResult, error) {
					return &pipeline.CheckMResult{
						Completeness: "99.5", Contamination: "8.0",
						GenomeSize: "5000000", N50: "100000",
					}, nil
				},
				RunAbricateFunc: func(_ context.Context, threads int,
					db, input, outputFile string) error {
					return writeAbricateOutput(outputFile)
				},
			}

			svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
				&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
			err := svc.Run(ctx, mock.ID)

			assert.NoError(t, err)
			assert.NotNil(t, updated)
			assert.Equal(t, models.AnalysisStatusDone, updated.Status)

			var results models.AnalysisResults
			assert.NoError(t, json.Unmarshal(updated.Metrics, &results))
			assert.NotEmpty(t, results.SecondarySpeciesName)
			assert.Equal(t, "Klebsiella pneumoniae",
				results.SecondarySpeciesName)
		})

	t.Run("Success - Secondary Species Not Written When Contamination <= 5%",
		func(t *testing.T) {
			mock := testmodels.CreateMockAnalysis()
			mock.Type = models.AnalysisTypeGenome
			mock.Status = models.AnalysisStatusPending
			fasta := filepath.Join(t.TempDir(), "contigs.fasta")
			os.WriteFile(fasta, []byte(">seq1\nATCG\n"), 0644)
			mock.Sample.Fastq1 = nil
			mock.Sample.Fastq2 = nil
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
				Config: pipeline.ToolsConfig{
					ResfinderDBPath: newResfinderRef(t),
				},
				RunCheckMFunc: func(_ context.Context, threads int,
					sample, assemblyDir, outputDir string) (
					*pipeline.CheckMResult, error) {
					return &pipeline.CheckMResult{
						Completeness: "99.5", Contamination: "0.5",
						GenomeSize: "5000000", N50: "100000",
					}, nil
				},
				RunAbricateFunc: func(_ context.Context, threads int,
					db, input, outputFile string) error {
					return writeAbricateOutput(outputFile)
				},
			}

			svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
				&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
			err := svc.Run(ctx, mock.ID)

			assert.NoError(t, err)
			assert.NotNil(t, updated)
			assert.Equal(t, models.AnalysisStatusDone, updated.Status)

			var results models.AnalysisResults
			assert.NoError(t, json.Unmarshal(updated.Metrics, &results))
			assert.Empty(t, results.SecondarySpeciesName)
		})
}

func TestAnalysisRunnerComplete(t *testing.T) {
	ctx := context.Background()

	originalConcurrency := config.AnalysisConcurrency
	config.AnalysisConcurrency = 4
	t.Cleanup(func() { config.AnalysisConcurrency = originalConcurrency })

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
			Config: pipeline.ToolsConfig{
				ResfinderDBPath: newResfinderRef(t),
			},
			RunFastQCFunc: func(_ context.Context, read1, read2,
				outputDir string) (string, string, error) {
				fastqcCalled = true
				return "qc1.html", "qc2.html", nil
			},
			RunAbricateFunc: func(_ context.Context, threads int,
				db, input, outputFile string) error {
				return writeAbricateOutput(outputFile)
			},
		}

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
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

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, zap.NewNop(), t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, pipeline.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			pipeline.ErrFastQC.Error())
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
			&mocks.MockCabgenPipeline{}, &mocks.MockCommander{}, &mocks.MockTaskEnqueuer{},
			zap.NewNop(), root)
		err := svc.Run(context.Background(), mock.ID)

		assert.NoError(t, err)
		for _, sub := range []string{"qc", "assembly", "amr", "report"} {
			assert.DirExists(t,
				filepath.Join(root, "uploads", "users",
					mock.UserID.String(), "samples",
					mock.SampleID.String(), "analyses",
					mock.ID.String(), sub))
		}
	})
}

func TestAnalysisRunnerZipResults(t *testing.T) {
	ctx := context.Background()

	newRepo := func(updated **models.Analysis, mock models.Analysis) *mocks.MockAnalysisRepository {
		return &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				analysis *models.Analysis) error {
				*updated = analysis
				return nil
			},
		}
	}

	t.Run("Success - Creates Zip and Sets ResultsZipPath", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		mock.ResultsZipPath = nil
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		updated := (*models.Analysis)(nil)
		repo := newRepo(&updated, mock)
		pl := &mocks.MockCabgenPipeline{
			RunFastQCFunc: func(_ context.Context, read1, read2,
				outputDir string) (string, string, error) {
				p1 := filepath.Join(outputDir, "reads1_fastqc.html")
				p2 := filepath.Join(outputDir, "reads2_fastqc.html")
				if err := os.WriteFile(p1, []byte("<html>"), 0644); err != nil {
					return "", "", err
				}
				if err := os.WriteFile(p2, []byte("<html>"), 0644); err != nil {
					return "", "", err
				}
				return p1, p2, nil
			},
		}
		enqueuer := &mocks.MockTaskEnqueuer{
			EnqueueContextFunc: func(_ context.Context, task *asynq.Task,
				_ ...asynq.Option) (*asynq.TaskInfo, error) {
				return &asynq.TaskInfo{ID: "t1", Queue: "emails"}, nil
			},
		}
		rootDir := t.TempDir()

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{}, enqueuer,
			zap.NewNop(), rootDir)
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusDone, updated.Status)
		assert.NotNil(t, updated.ResultsZipPath)

		expectedZip := filepath.Join(rootDir, "uploads", "users",
			mock.UserID.String(), "samples", mock.SampleID.String(),
			"analyses", mock.ID.String(), "report",
			mock.Sample.OriginCode+"_FASTQC_results.zip")
		assert.Equal(t, expectedZip, *updated.ResultsZipPath)
		assert.FileExists(t, expectedZip)

		zr, err := zip.OpenReader(expectedZip)
		assert.NoError(t, err)
		defer zr.Close()

		var names []string
		for _, f := range zr.File {
			names = append(names, f.Name)
		}
		assert.Contains(t, names, filepath.Join(mock.ID.String(),
			"qc", "reads1_fastqc.html"))
		assert.Contains(t, names, filepath.Join(mock.ID.String(),
			"qc", "reads2_fastqc.html"))
		assert.NotContains(t, names, filepath.Join(mock.ID.String(),
			"report", mock.Sample.OriginCode+"_FASTQC_results.zip"))
	})

	t.Run("Warning - Zip Failure Does Not Fail Analysis", func(t *testing.T) {
		mock := testmodels.CreateMockAnalysis()
		mock.Type = models.AnalysisTypeFastQC
		mock.Status = models.AnalysisStatusPending
		mock.ResultsZipPath = nil
		fq1, fq2 := "r1.fq", "r2.fq"
		mock.Sample.Fastq1 = &fq1
		mock.Sample.Fastq2 = &fq2

		updated := (*models.Analysis)(nil)
		repo := newRepo(&updated, mock)
		pl := &mocks.MockCabgenPipeline{
			RunFastQCFunc: func(_ context.Context, read1, read2,
				outputDir string) (string, string, error) {
				dir := filepath.Dir(outputDir)
				if err := os.RemoveAll(dir); err != nil {
					return "", "", err
				}
				return "", "", nil
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zap.WarnLevel)

		svc := services.NewAnalysisRunnerService(repo, pl, &mocks.MockCommander{},
			&mocks.MockTaskEnqueuer{}, mockLogger, t.TempDir())
		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusDone, updated.Status)
		assert.Nil(t, updated.ResultsZipPath)
		assert.GreaterOrEqual(t, logs.Len(), 1)
	})
}
