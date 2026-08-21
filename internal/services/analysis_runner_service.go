package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const SecondarySpeciesContaminationThreshold = 5.0

var (
	versionsOnce sync.Once
	versions     []pipeline.ToolVersion
)

type AnalysisRunnerFolders struct {
	QCDir       string
	AssemblyDir string
	AMRDir      string
	ReportDir   string
}

type AnalysisRunnerService interface {
	Run(ctx context.Context, analysisID uuid.UUID) error
}

type analysisRunnerService struct {
	Repo        repositories.AnalysisRepository
	Pipeline    pipeline.CabgenPipeline
	Commander   pipeline.Commander
	AsynqClient TaskEnqueuer
	Logger      *zap.Logger
	RootDir     string
}

func NewAnalysisRunnerService(
	repo repositories.AnalysisRepository,
	pipeline pipeline.CabgenPipeline,
	commander pipeline.Commander,
	asynqClient TaskEnqueuer,
	logger *zap.Logger, rootDir string) AnalysisRunnerService {
	return &analysisRunnerService{
		Repo:        repo,
		Pipeline:    pipeline,
		Commander:   commander,
		AsynqClient: asynqClient,
		Logger:      logger,
		RootDir:     rootDir,
	}
}

func (s *analysisRunnerService) prepareFolders(
	userID, sampleID, analysisID string) (*AnalysisRunnerFolders, error) {
	rootDir := filepath.Join(s.RootDir, "uploads", "users", userID,
		"samples", sampleID, "analyses", analysisID)

	qcDir := filepath.Join(rootDir, "qc")
	assemblyDir := filepath.Join(rootDir, "assembly")
	amrDir := filepath.Join(rootDir, "amr")
	reportDir := filepath.Join(rootDir, "report")

	for _, dir := range []string{qcDir, assemblyDir, amrDir, reportDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AnalysisRunnerService", "prepareFolders",
				logging.CreateFolderError, err,
			)...)
			return nil, ErrCreateFolder
		}
	}

	return &AnalysisRunnerFolders{
		QCDir: qcDir, AssemblyDir: assemblyDir, AMRDir: amrDir,
		ReportDir: reportDir,
	}, nil
}

func (s *analysisRunnerService) getVersions(ctx context.Context,
) {
	versionsOnce.Do(func() {
		versions = pipeline.GetBioinfoProgramVersions(ctx, s.Commander)
	})
}

func (s *analysisRunnerService) runFastQC(ctx context.Context,
	analysis *models.Analysis, outputDir string) error {
	s.updateStep(ctx, analysis, models.StepFastQC)
	s.Logger.Info(
		fmt.Sprintf("%s: Started FastQC step", analysis.ID.String()),
		logging.ServiceInfoLogging("AnalysisRunnerService", "runFastQC",
			"CabgenPipeline")...,
	)

	fastq1Path, ok := utils.ResolveSampleFilePath(
		s.RootDir,
		analysis.UserID.String(),
		analysis.SampleID.String(),
		*analysis.Sample.Fastq1,
		"fastq", "",
	)
	if !ok {
		return fmt.Errorf("fastq1 file not found: %s", *analysis.Sample.Fastq1)
	}
	fastq2Path, ok := utils.ResolveSampleFilePath(
		s.RootDir,
		analysis.UserID.String(),
		analysis.SampleID.String(),
		*analysis.Sample.Fastq2,
		"fastq", "",
	)
	if !ok {
		return fmt.Errorf("fastq2 file not found: %s", *analysis.Sample.Fastq2)
	}

	fastqc1, fastqc2, err := s.Pipeline.RunFastQC(
		ctx, fastq1Path, fastq2Path, outputDir)
	if err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed FastQC step: %v", analysis.ID.String(), err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runFastQC",
				logging.AnalysisRunError, err,
			)...)
		if isInputError(err) {
			return err
		}
		return pipeline.ErrFastQC
	}

	analysis.FastQC1 = &fastqc1
	analysis.FastQC2 = &fastqc2
	if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed to update analysis in FastQC step: %v",
			analysis.ID.String(), err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runFastQC",
				logging.DatabaseError, err,
			)...)
		return ErrInternal
	}

	return nil
}

func (s *analysisRunnerService) runGenome(ctx context.Context,
	analysis *models.Analysis, results *models.AnalysisResults,
	folders *AnalysisRunnerFolders) error {
	s.Logger.Info(
		fmt.Sprintf("%s: Started Genome step", analysis.ID.String()),
		logging.ServiceInfoLogging("AnalysisRunnerService", "runGenome",
			"CabgenPipeline")...,
	)

	// Using 80% of total cores
	threads := int(math.Round(
		(float64(runtime.NumCPU()) * 0.8) /
			float64(config.AnalysisConcurrency)))
	var assemblyPath *string

	if analysis.Sample.Fasta != nil {
		resolved, ok := utils.ResolveSampleFilePath(
			s.RootDir,
			analysis.UserID.String(),
			analysis.SampleID.String(),
			*analysis.Sample.Fasta,
			"fasta", analysis.ID.String(),
		)
		if !ok {
			s.Logger.Warn(fmt.Sprintf(
				"%s: FASTA file not found at %s, falling back to reads",
				analysis.ID.String(), *analysis.Sample.Fasta),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.MissingFileError, fmt.Errorf("file not found"),
				)...)
			analysis.Sample.Fasta = nil
		} else {
			dst := filepath.Join(folders.AssemblyDir,
				filepath.Base(*analysis.Sample.Fasta))
			if resolved != dst {
				if err := utils.CopyFile(resolved, dst); err != nil {
					return fmt.Errorf("failed to prepare assembly: %w", err)
				}
			}
			assemblyPath = &dst
		}
	}

	if analysis.Sample.Fastq1 != nil && analysis.Sample.Fastq2 != nil &&
		analysis.Sample.Fasta == nil {
		s.updateStep(ctx, analysis, models.StepUnicycler)

		fastq1Path, ok := utils.ResolveSampleFilePath(
			s.RootDir,
			analysis.UserID.String(),
			analysis.SampleID.String(),
			*analysis.Sample.Fastq1,
			"fastq", "",
		)
		if !ok {
			return fmt.Errorf("fastq1 file not found: %s", *analysis.Sample.Fastq1)
		}
		fastq2Path, ok := utils.ResolveSampleFilePath(
			s.RootDir,
			analysis.UserID.String(),
			analysis.SampleID.String(),
			*analysis.Sample.Fastq2,
			"fastq", "",
		)
		if !ok {
			return fmt.Errorf("fastq2 file not found: %s", *analysis.Sample.Fastq2)
		}

		assemblyOutPath := fmt.Sprintf("%s_assembly.fasta",
			analysis.Sample.OriginCode)
		assembly, err := s.Pipeline.RunUnicycler(ctx, threads,
			fastq1Path, fastq2Path,
			s.Pipeline.GetConfig().SpadesPath, folders.AssemblyDir,
			assemblyOutPath)
		if err != nil {
			s.Logger.Error(fmt.Sprintf(
				"%s: Failed Genome step - Unicycler: %v",
				analysis.ID.String(), err),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
			if isInputError(err) {
				return err
			}
			return pipeline.ErrUnicycler
		}
		assemblyPath = &assembly
		assemblyFileName := filepath.Base(assembly)
		analysis.Sample.Fasta = &assemblyFileName

		if err := s.Repo.UpdateSample(ctx, &analysis.Sample); err != nil {
			s.Logger.Warn(fmt.Sprintf(
				"%s: Failed to persist assembly path to sample",
				analysis.ID.String()),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
		}
	}

	if assemblyPath == nil {
		return fmt.Errorf("no input files: need FASTA or FASTQ pair")
	}

	prokkaOutDir := filepath.Join(folders.AssemblyDir, "prokka")
	s.updateStep(ctx, analysis, models.StepProkka)
	if err := s.Pipeline.RunProkka(ctx, threads, *assemblyPath,
		prokkaOutDir); err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed Genome step - Prokka: %v", analysis.ID.String(),
			err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runGenome",
				logging.AnalysisRunError, err,
			)...)
		if isInputError(err) {
			return err
		}
		return pipeline.ErrProkka
	}

	ext := filepath.Ext(*assemblyPath)
	checkmSample := strings.TrimSuffix(filepath.Base(*assemblyPath), ext)
	checkMOutput := filepath.Join(folders.AssemblyDir, "checkm_output")
	s.updateStep(ctx, analysis, models.StepCheckM)
	checkmResult, err := s.Pipeline.RunCheckM(ctx, threads, checkmSample,
		folders.AssemblyDir, checkMOutput)
	if err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed Genome step - CheckM: %v", analysis.ID.String(),
			err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runGenome",
				logging.AnalysisRunError, err,
			)...)
		return pipeline.ErrCheckM
	}

	if checkmResult != nil {
		results.CheckMCompleteness = checkmResult.Completeness
		results.CheckMContamination = checkmResult.Contamination
		results.CheckMGenomeSize = checkmResult.GenomeSize
		results.CheckMN50 = checkmResult.N50
	}

	s.updateStep(ctx, analysis, models.StepKraken2)
	krakenResult1, krakenResult2, err := s.Pipeline.RunKraken2(ctx, threads,
		*assemblyPath, folders.AssemblyDir)
	if err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed Genome step - Kraken2: %v", analysis.ID.String(),
			err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runGenome",
				logging.AnalysisRunError, err,
			)...)
		return pipeline.ErrKraken2
	}

	if krakenResult1 != nil {
		s.updateStep(ctx, analysis, models.StepSpecies)
		speciesResult, err := s.Pipeline.ProcessSpecies(ctx, threads,
			analysis.SampleID.String(), krakenResult1.Name, *assemblyPath,
			folders.AssemblyDir)
		if err != nil {
			s.Logger.Error(fmt.Sprintf(
				"%s: Failed Genome step - Species: %v", analysis.ID.String(),
				err),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
			return pipeline.ErrSpecies
		}

		if speciesResult != nil {
			results.PrimarySpeciesName = speciesResult.DisplayName
			results.MLST = speciesResult.MLSTSpecies
			results.PoliMutations = speciesResult.PoliMutations
			results.OtherMutations = speciesResult.OtherMutations
		}
	}

	if krakenResult2 != nil {
		contamination, err := strconv.ParseFloat(
			results.CheckMContamination, 32)
		if err != nil {
			s.Logger.Warn(fmt.Sprintf(
				"%s: Invalid CheckM contamination value: %q",
				analysis.ID.String(), results.CheckMContamination),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
		}
		if contamination > SecondarySpeciesContaminationThreshold {
			results.SecondarySpeciesName = krakenResult2.Name
		}
	}

	abricateInput := filepath.Join(prokkaOutDir, "genome.ffn")
	abricateDBs := map[string]string{
		"resfinder": filepath.Join(folders.AMRDir, fmt.Sprintf(
			"%s_outAbricateRes", analysis.SampleID.String())),
		"vfdb": filepath.Join(folders.AMRDir, fmt.Sprintf(
			"%s_outAbricateVFDB", analysis.SampleID.String())),
		"plasmidfinder": filepath.Join(folders.AMRDir, fmt.Sprintf(
			"%s_outAbricatePlasmid", analysis.SampleID.String())),
	}
	s.updateStep(ctx, analysis, models.StepAbricate)
	for db, outputFile := range abricateDBs {
		if err := s.Pipeline.RunAbricate(ctx, threads, db, abricateInput,
			outputFile); err != nil {
			s.Logger.Error(fmt.Sprintf(
				"%s: Failed Genome step - Abricate (%s): %v",
				analysis.ID.String(), db, err),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
			return pipeline.ErrAbricate
		}

		rawResult, err := pipeline.GetAbricateResult(outputFile)
		if err != nil {
			s.Logger.Error(fmt.Sprintf(
				"%s: Failed Genome step - Abricate Result (%s): %v",
				analysis.ID.String(), db, err),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
			return pipeline.ErrAbricate
		}

		switch db {
		case "resfinder":
			genes, err := pipeline.ProcessResfinder(rawResult,
				s.Pipeline.GetConfig().ResfinderDBPath)
			if err != nil {
				s.Logger.Error(fmt.Sprintf(
					"%s: Failed Genome step - ProcessResfinder: %v",
					analysis.ID.String(), err),
					logging.ServiceLogging(
						"AnalysisRunnerService", "runGenome",
						logging.AnalysisRunError, err,
					)...)
				return pipeline.ErrAbricate
			}
			results.AcquiredResistance = genes
		case "vfdb":
			results.VFDB = pipeline.ProcessVFDB(rawResult)
		case "plasmidfinder":
			results.PlasmidFinder = pipeline.ProcessPlasmidFinder(rawResult)
		}
	}

	genomeSize, _ := strconv.Atoi(results.CheckMGenomeSize)
	if analysis.Sample.Fastq1 != nil && analysis.Sample.Fastq2 != nil &&
		genomeSize > 0 {
		s.updateStep(ctx, analysis, models.StepCoverage)
		fastq1Path, ok := utils.ResolveSampleFilePath(
			s.RootDir,
			analysis.UserID.String(),
			analysis.SampleID.String(),
			*analysis.Sample.Fastq1,
			"fastq", "",
		)
		if !ok {
			return fmt.Errorf("fastq1 file not found: %s", *analysis.Sample.Fastq1)
		}
		fastq2Path, ok := utils.ResolveSampleFilePath(
			s.RootDir,
			analysis.UserID.String(),
			analysis.SampleID.String(),
			*analysis.Sample.Fastq2,
			"fastq", "",
		)
		if !ok {
			return fmt.Errorf("fastq2 file not found: %s", *analysis.Sample.Fastq2)
		}
		coverage, err := pipeline.CalculateCoverage(
			fastq1Path, fastq2Path, int64(genomeSize))
		if err != nil {
			s.Logger.Error(fmt.Sprintf(
				"%s: Failed Genome step - CalculateCoverage: %v",
				analysis.ID.String(), err),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
		} else {
			results.Coverage = coverage
		}
	}

	return nil
}

func (s *analysisRunnerService) runComplete(ctx context.Context,
	analysis *models.Analysis, results *models.AnalysisResults,
	folders *AnalysisRunnerFolders) error {

	if err := s.runFastQC(ctx, analysis, folders.QCDir); err != nil {
		return err
	}

	if err := s.runGenome(ctx, analysis, results, folders); err != nil {
		return err
	}

	return nil
}

func (s *analysisRunnerService) updateStep(ctx context.Context,
	analysis *models.Analysis, step models.AnalysisStep) {
	analysis.Step = step
	if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
		s.Logger.Warn("Service Warning", logging.ServiceLogging(
			"AnalysisRunnerService", "updateStep",
			logging.DatabaseError, err,
		)...)
	}
}

func (s *analysisRunnerService) finalizeAnalysis(ctx context.Context,
	analysis *models.Analysis, results *models.AnalysisResults, runErr error) {
	finished := time.Now()
	analysis.FinishedAt = &finished
	analysis.Step = ""

	results.Versions = versions

	if runErr != nil {
		analysis.Status = models.AnalysisStatusFailed
		msg := runErr.Error()
		analysis.ErrorMessage = &msg
	} else {
		analysis.Status = models.AnalysisStatusDone
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		s.Logger.Warn(fmt.Sprintf(
			"%s: Failed to marshal analysis results: %v",
			analysis.ID.String(), err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "finalizeAnalysis",
				logging.AnalysisRunError, err,
			)...)
	} else {
		analysis.Metrics = jsonData
	}

	if analysis.Status == models.AnalysisStatusDone {
		s.zipAnalysisResults(analysis)
	}

	if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisRunnerService", "Run",
			logging.DatabaseError, err,
		)...)
	}
}

func (s *analysisRunnerService) zipAnalysisResults(
	analysis *models.Analysis) {
	analysisFolder := filepath.Join(s.RootDir, "uploads", "users",
		analysis.UserID.String(), "samples", analysis.SampleID.String(),
		"analyses", analysis.ID.String())

	reportDir := filepath.Join(analysisFolder, "report")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		s.Logger.Warn("Service Warning", logging.ServiceLogging(
			"AnalysisRunnerService", "zipAnalysisResults",
			logging.CreateFolderError, err,
		)...)
		return
	}

	zipName := utils.SanitizeFilename(analysis.Sample.OriginCode) + "_" +
		string(analysis.Type) + "_results.zip"
	zipPath := filepath.Join(reportDir, zipName)
	if err := utils.ZipSubdirectories(analysisFolder,
		[]string{"qc", "assembly", "amr"}, zipPath); err != nil {
		s.Logger.Warn("Service Warning", logging.ServiceLogging(
			"AnalysisRunnerService", "zipAnalysisResults",
			logging.CreateFolderError, err,
		)...)
		return
	}

	analysis.ResultsZipPath = &zipPath
}

func (s *analysisRunnerService) Run(ctx context.Context,
	analysisID uuid.UUID) error {
	analysis, err := s.Repo.GetAnalysisByID(ctx, analysisID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisRunnerService", "Run",
			logging.DatabaseNotFoundError, err,
		)...)
		return ErrNotFound
	}
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisRunnerService", "Run",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	s.getVersions(ctx)

	start := time.Now()
	analysis.Status = models.AnalysisStatusRunning
	analysis.StartedAt = &start
	if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisRunnerService", "Run",
			logging.DatabaseError, err,
		)...)
		return ErrInternal
	}

	var results models.AnalysisResults

	folders, err := s.prepareFolders(analysis.UserID.String(),
		analysis.SampleID.String(), analysis.ID.String())
	if err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed to prepare folders: %v", analysisID.String(), err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "Run",
				logging.CreateFolderError, err,
			)...)
		s.finalizeAnalysis(ctx, analysis, &results, pipeline.ErrPrepareFolders)
		return pipeline.ErrAnalysisRun
	}

	s.Logger.Info(
		fmt.Sprintf("Analysis %s started (type: %s)", analysisID.String(),
			analysis.Type),
		logging.ServiceInfoLogging("AnalysisRunnerService", "Run",
			"CabgenPipeline")...,
	)

	var runErr error
	switch analysis.Type {
	case models.AnalysisTypeFastQC:
		runErr = s.runFastQC(ctx, analysis, folders.QCDir)
	case models.AnalysisTypeGenome:
		runErr = s.runGenome(ctx, analysis, &results, folders)
	case models.AnalysisTypeComplete:
		runErr = s.runComplete(ctx, analysis, &results, folders)
	default:
		s.Logger.Error(fmt.Sprintf(
			"Analysis %s: unknown analysis type %s", analysisID.String(),
			analysis.Type),
			logging.ServiceLogging(
				"AnalysisRunnerService", "Run",
				logging.AnalysisRunError,
				fmt.Errorf("unknown type: %s", analysis.Type),
			)...)
		runErr = pipeline.ErrUnknownAnalysisType
	}

	s.finalizeAnalysis(ctx, analysis, &results, runErr)

	shouldEnqueueEmail := runErr == nil
	if !shouldEnqueueEmail {
		if retryCount, ok := asynq.GetRetryCount(ctx); ok {
			if maxRetry, ok := asynq.GetMaxRetry(ctx); ok {
				shouldEnqueueEmail = retryCount >= maxRetry
			}
		}
	}

	if shouldEnqueueEmail {
		task, err := tasks.NewAnalysisDoneEmailTask(analysisID)
		if err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"AnalysisRunnerService", "Run", logging.AsynqTaskError, err,
			)...)
		} else {
			info, err := s.AsynqClient.EnqueueContext(ctx, task,
				asynq.Queue(tasks.QueueEmail))
			if err != nil {
				s.Logger.Error("Service Error", logging.ServiceLogging(
					"AnalysisRunnerService", "Run",
					logging.RedisDispatchError, err,
				)...)
			} else {
				s.Logger.Info("Redis Task Info",
					logging.ServiceInfoLogging(
						"AnalysisRunnerService", "Run",
						logging.TaskEnqueuedSuccess,
						zap.String("task_id", info.ID),
						zap.String("queue", info.Queue),
					)...)
			}
		}
	}

	if runErr != nil {
		s.Logger.Error(fmt.Sprintf(
			"Analysis %s failed after %v: %v", analysisID.String(),
			time.Since(start), runErr),
			logging.ServiceLogging(
				"AnalysisRunnerService", "Run",
				logging.AnalysisRunError, runErr,
			)...)
		return pipeline.ErrAnalysisRun
	}

	s.Logger.Info(
		fmt.Sprintf("Analysis %s completed in %v", analysisID.String(),
			time.Since(start)),
		logging.ServiceInfoLogging("AnalysisRunnerService", "Run",
			"CabgenPipeline")...,
	)

	return nil
}

func isInputError(err error) bool {
	return errors.Is(err, pipeline.ErrCorruptedInput) ||
		errors.Is(err, pipeline.ErrEmptyReads) ||
		errors.Is(err, pipeline.ErrInvalidFormat) ||
		errors.Is(err, pipeline.ErrFileNotFound)
}
