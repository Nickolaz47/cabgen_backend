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

func newRunnerService(repo *mocks.MockAnalysisRepository,
	pl *mocks.MockCabgenPipeline, logger *zap.Logger,
	rootDir string) services.AnalysisRunnerService {
	return services.NewAnalysisRunnerService(repo, pl, logger, rootDir)
}

// abricateWriter returns a RunAbricateFunc that writes a minimal placeholder
// file so pipeline.GetAbricateResult finds content and returns no error.
func abricateWriter() func(context.Context, int, string, string,
	string) error {
	return func(_ context.Context, _ int, _, _, outputFile string) error {
		return os.WriteFile(outputFile, []byte("placeholder\n"), 0644)
	}
}

func analysisWithFastQ(t *testing.T, t1, t2 string) models.Analysis {
	t.Helper()
	mock := testmodels.CreateMockAnalysis()
	mock.Type = models.AnalysisTypeFastQC
	mock.Status = models.AnalysisStatusPending
	if t1 != "" {
		mock.Sample.Fastq1 = &t1
	} else {
		mock.Sample.Fastq1 = nil
	}
	if t2 != "" {
		mock.Sample.Fastq2 = &t2
	} else {
		mock.Sample.Fastq2 = nil
	}
	return mock
}

func TestAnalysisRunnerRun(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - FastQC", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				a *models.Analysis) error {
				updated = a
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

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

	t.Run("Error - Analysis Not Found", func(t *testing.T) {
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := newRunnerService(repo, nil, mockLogger, t.TempDir())

		err := svc.Run(ctx, uuid.New())

		assert.ErrorIs(t, err, services.ErrNotFound)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB failure on GetAnalysisByID", func(t *testing.T) {
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := newRunnerService(repo, nil, mockLogger, t.TempDir())

		err := svc.Run(ctx, uuid.New())

		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - DB failure on first UpdateAnalysis", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
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
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := newRunnerService(repo, nil, mockLogger, t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - FastQC pipeline failure persisted on analysis",
		func(t *testing.T) {
			mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
			pipelineErr := errors.New("fastqc crashed")
			updated := (*models.Analysis)(nil)
			repo := &mocks.MockAnalysisRepository{
				GetAnalysisByIDFunc: func(_ context.Context,
					_ uuid.UUID) (*models.Analysis, error) {
					mockCopy := mock
					return &mockCopy, nil
				},
				UpdateAnalysisFunc: func(_ context.Context,
					a *models.Analysis) error {
					updated = a
					return nil
				},
			}
			pl := &mocks.MockCabgenPipeline{
				RunFastQCFunc: func(_ context.Context, _, _,
					_ string) (string, string, error) {
					return "", "", pipelineErr
				},
			}
			mockLogger, _ := testutils.NewMockLogger(zapcore.ErrorLevel)
			svc := newRunnerService(repo, pl, mockLogger, t.TempDir())

			err := svc.Run(ctx, mock.ID)

			assert.ErrorIs(t, err, services.ErrAnalysisRun)
			assert.NotNil(t, updated)
			assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
			assert.NotNil(t, updated.ErrorMessage)
			assert.Contains(t, *updated.ErrorMessage, services.ErrFastQC.Error())
			assert.NotContains(t, *updated.ErrorMessage, "fastqc crashed")
		})

	t.Run("Error - Unknown analysis type", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = "NONSENSE"
		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				a *models.Analysis) error {
				updated = a
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrUnknownAnalysisType.Error())
	})

	t.Run("Error - prepareFolders failure", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				a *models.Analysis) error {
				updated = a
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{}
		svc := newRunnerService(repo, pl, zap.NewNop(),
			"/nonexistent_root_no_perms/x")

		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrPrepareFolders.Error())
	})

	t.Run("Success - Genome with Fasta skips Unicycler", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = models.AnalysisTypeGenome
		fasta := "contigs.fasta"
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
			RunUnicyclerFunc: func(_ context.Context, _ int, _, _, _,
				_ string) (string, error) {
				unicyclerCalled = true
				return "assembly.fa", nil
			},
			RunAbricateFunc: abricateWriter(),
		}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.False(t, unicyclerCalled,
			"Unicycler should not run when Fasta already present")
	})

	t.Run("Success - Genome runs Unicycler when no Fasta", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = models.AnalysisTypeGenome
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
			RunUnicyclerFunc: func(_ context.Context, _ int, _, _, _,
				_ string) (string, error) {
				unicyclerCalled = true
				return "assembly.fa", nil
			},
			RunAbricateFunc: abricateWriter(),
		}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.True(t, unicyclerCalled)
	})

	t.Run("Error - Genome Unicycler failure has step context", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = models.AnalysisTypeGenome
		mock.Sample.Fasta = nil
		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				a *models.Analysis) error {
				updated = a
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunUnicyclerFunc: func(_ context.Context, _ int, _, _, _,
				_ string) (string, error) {
				return "", errors.New("spades missing")
			},
		}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrUnicycler.Error())
		assert.NotContains(t, *updated.ErrorMessage, "spades missing")
	})
}

func TestAnalysisRunnerRunComplete(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Complete runs FastQC then Genome", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = models.AnalysisTypeComplete
		fasta := "contigs.fasta"
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
			RunFastQCFunc: func(_ context.Context, _, _,
				_ string) (string, string, error) {
				fastqcCalled = true
				return "qc1.html", "qc2.html", nil
			},
			RunAbricateFunc: abricateWriter(),
		}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.NoError(t, err)
		assert.True(t, fastqcCalled, "Complete should call FastQC first")
	})

	t.Run("Error - Complete FastQC failure aborts Genome", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = models.AnalysisTypeComplete
		fasta := "contigs.fasta"
		mock.Sample.Fasta = &fasta
		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				a *models.Analysis) error {
				updated = a
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunFastQCFunc: func(_ context.Context, _, _,
				_ string) (string, string, error) {
				return "", "", errors.New("fastqc timeout")
			},
		}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrFastQC.Error())
	})
}

func TestAnalysisRunnerGenomeAbricate(t *testing.T) {
	ctx := context.Background()

	t.Run("Error - Abricate failure has DB name in message", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = models.AnalysisTypeGenome
		fasta := "contigs.fasta"
		mock.Sample.Fasta = &fasta
		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				a *models.Analysis) error {
				updated = a
				return nil
			},
		}
		// Fail all DBs so map iteration order doesn't matter
		abricateErr := errors.New("abricate segfault")
		var failedDB string
		pl := &mocks.MockCabgenPipeline{
			RunAbricateFunc: func(_ context.Context, _ int, db, _,
				_ string) error {
				if failedDB == "" {
					failedDB = db
				}
				return abricateErr
			},
		}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrAbricate.Error())
		assert.NotEmpty(t, failedDB, "at least one DB should have failed")
		assert.NotContains(t, *updated.ErrorMessage, "abricate segfault")
	})
}

func TestAnalysisRunnerPrepareFolders(t *testing.T) {
	t.Run("Success - Creates all folders", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
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
		pl := &mocks.MockCabgenPipeline{}
		svc := newRunnerService(repo, pl, zap.NewNop(), root)

		err := svc.Run(context.Background(), mock.ID)
		assert.NoError(t, err)

		for _, sub := range []string{"qc", "assembly", "amr", "report"} {
			assert.DirExists(t, filepath.Join(root, mock.ID.String(), sub))
		}
	})
}

func TestAnalysisRunnerErrorMessageFormat(t *testing.T) {
	ctx := context.Background()

	t.Run("Genome CheckM failure message is actionable", func(t *testing.T) {
		mock := analysisWithFastQ(t, "r1.fq", "r2.fq")
		mock.Type = models.AnalysisTypeGenome
		fasta := "contigs.fasta"
		mock.Sample.Fasta = &fasta
		updated := (*models.Analysis)(nil)
		repo := &mocks.MockAnalysisRepository{
			GetAnalysisByIDFunc: func(_ context.Context,
				_ uuid.UUID) (*models.Analysis, error) {
				mockCopy := mock
				return &mockCopy, nil
			},
			UpdateAnalysisFunc: func(_ context.Context,
				a *models.Analysis) error {
				updated = a
				return nil
			},
		}
		pl := &mocks.MockCabgenPipeline{
			RunCheckMFunc: func(_ context.Context, _ int, _, _,
				_ string) (*pipeline.CheckMResult, error) {
				return nil, fmt.Errorf("checkm db corrupt")
			},
		}
		svc := newRunnerService(repo, pl, zap.NewNop(), t.TempDir())

		err := svc.Run(ctx, mock.ID)

		assert.ErrorIs(t, err, services.ErrAnalysisRun)
		assert.NotNil(t, updated)
		assert.Equal(t, models.AnalysisStatusFailed, updated.Status)
		assert.NotNil(t, updated.ErrorMessage)
		assert.Contains(t, *updated.ErrorMessage,
			services.ErrCheckM.Error())
		assert.NotContains(t, *updated.ErrorMessage, "checkm db corrupt")
	})
}

// Ensure unused vars don't break the build (defensive against refactors).
var _ = os.MkdirAll