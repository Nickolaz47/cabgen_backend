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
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	Repo     repositories.AnalysisRepository
	Pipeline pipeline.CabgenPipeline
	Logger   *zap.Logger
	RootDir  string
}

func NewAnalysisRunnerService(
	repo repositories.AnalysisRepository,
	pipeline pipeline.CabgenPipeline,
	logger *zap.Logger, rootDir string) AnalysisRunnerService {
	return &analysisRunnerService{
		Repo:     repo,
		Pipeline: pipeline,
		Logger:   logger,
		RootDir:  rootDir,
	}
}

func (s *analysisRunnerService) prepareFolders(
	analysisID string) (*AnalysisRunnerFolders, error) {
	rootDir := filepath.Join(s.RootDir, analysisID)

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

func (s *analysisRunnerService) runFastQC(ctx context.Context,
	analysis *models.Analysis, outputDir string) error {
	s.Logger.Info(
		fmt.Sprintf("%s: Started FastQC step", analysis.ID.String()),
		logging.ServiceInfoLogging("AnalysisRunnerService", "runFastQC",
			"CabgenPipeline")...,
	)

	fastqc1, fastqc2, err := s.Pipeline.RunFastQC(
		ctx, *analysis.Sample.Fastq1, *analysis.Sample.Fastq2, outputDir)
	if err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed FastQC step: %v", analysis.ID.String(), err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runFastQC",
				logging.AnalysisRunError, err,
			)...)
		return ErrFastQC
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

	// Using 20% of the total cores
	threads := int(math.Round((float64(runtime.NumCPU()) * 0.8) / 4))
	assemblyPath := analysis.Sample.Fasta

	if analysis.Sample.Fastq1 != nil && analysis.Sample.Fastq2 != nil &&
		analysis.Sample.Fasta == nil {
		assembly, err := s.Pipeline.RunUnicycler(ctx, threads,
			*analysis.Sample.Fastq1, *analysis.Sample.Fastq2,
			s.Pipeline.GetConfig().SpadesPath, folders.AssemblyDir)
		if err != nil {
			s.Logger.Error(fmt.Sprintf(
				"%s: Failed Genome step - Unicycler: %v",
				analysis.ID.String(), err),
				logging.ServiceLogging(
					"AnalysisRunnerService", "runGenome",
					logging.AnalysisRunError, err,
				)...)
			return ErrUnicycler
		}
		assemblyPath = &assembly
		analysis.Sample.Fasta = &assembly
	}

	prokkaOutDir := filepath.Join(folders.AssemblyDir, "prokka")
	if err := s.Pipeline.RunProkka(ctx, threads, *assemblyPath,
		prokkaOutDir); err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed Genome step - Prokka: %v", analysis.ID.String(),
			err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runGenome",
				logging.AnalysisRunError, err,
			)...)
		return ErrProkka
	}

	ext := filepath.Ext(*assemblyPath)
	checkmSample := strings.TrimSuffix(filepath.Base(*assemblyPath), ext)
	checkmResult, err := s.Pipeline.RunCheckM(ctx, threads, checkmSample,
		folders.AssemblyDir, folders.AssemblyDir)
	if err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed Genome step - CheckM: %v", analysis.ID.String(),
			err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "runGenome",
				logging.AnalysisRunError, err,
			)...)
		return ErrCheckM
	}

	if checkmResult != nil {
		results.CheckMCompleteness = checkmResult.Completeness
		results.CheckMContamination = checkmResult.Contamination
		results.CheckMGenomeSize = checkmResult.GenomeSize
		results.CheckMN50 = checkmResult.N50
	}

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
		return ErrKraken2
	}

	if krakenResult1 != nil {
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
			return ErrSpecies
		}

		if speciesResult != nil {
			results.PrimarySpeciesName = speciesResult.DisplayName
			results.MLST = speciesResult.MLSTSpecies
			results.PoliMutations = speciesResult.PoliMutations
			results.OtherMutations = speciesResult.OtherMutations
		}
	}

	if krakenResult2 != nil {
		results.SecondarySpeciesName = krakenResult2.Name
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
			return ErrAbricate
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
			return ErrAbricate
		}

		switch db {
		case "resfinder":
			genes, blast, err := pipeline.ProcessResfinder(rawResult,
				s.Pipeline.GetConfig().ResfinderDBPath)
			if err == nil {
				results.ResfinderGenes = genes
				results.ResfinderBlast = blast
			}
		case "vfdb":
			results.VFDB = pipeline.ProcessVFDB(rawResult)
		case "plasmidfinder":
			results.PlasmidFinder = pipeline.ProcessPlasmidFinder(rawResult)
		}
	}

	genomeSize, _ := strconv.Atoi(results.CheckMGenomeSize)
	if analysis.Sample.Fastq1 != nil && analysis.Sample.Fastq2 != nil &&
		genomeSize > 0 {
		coverage, err := pipeline.CalculateCoverage(
			*analysis.Sample.Fastq1, *analysis.Sample.Fastq2,
			int64(genomeSize))
		if err == nil {
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

func (s *analysisRunnerService) finalizeAnalysis(ctx context.Context,
	analysis *models.Analysis, results *models.AnalysisResults, runErr error) {
	finished := time.Now()
	analysis.FinishedAt = &finished

	if runErr != nil {
		analysis.Status = models.AnalysisStatusFailed
		msg := runErr.Error()
		analysis.ErrorMessage = &msg
	} else {
		analysis.Status = models.AnalysisStatusDone
	}

	jsonData, err := json.Marshal(results)
	if err == nil {
		analysis.Metrics = jsonData
	}

	if err := s.Repo.UpdateAnalysis(ctx, analysis); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"AnalysisRunnerService", "Run",
			logging.DatabaseError, err,
		)...)
	}
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

	folders, err := s.prepareFolders(analysis.ID.String())
	if err != nil {
		s.Logger.Error(fmt.Sprintf(
			"%s: Failed to prepare folders: %v", analysisID.String(), err),
			logging.ServiceLogging(
				"AnalysisRunnerService", "Run",
				logging.CreateFolderError, err,
			)...)
		s.finalizeAnalysis(ctx, analysis, &results, ErrPrepareFolders)
		return ErrAnalysisRun
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
		runErr = ErrUnknownAnalysisType
	}

	s.finalizeAnalysis(ctx, analysis, &results, runErr)

	if runErr != nil {
		s.Logger.Error(fmt.Sprintf(
			"Analysis %s failed after %v: %v", analysisID.String(),
			time.Since(start), runErr),
			logging.ServiceLogging(
				"AnalysisRunnerService", "Run",
				logging.AnalysisRunError, runErr,
			)...)
		return ErrAnalysisRun
	}

	s.Logger.Info(
		fmt.Sprintf("Analysis %s completed in %v", analysisID.String(),
			time.Since(start)),
		logging.ServiceInfoLogging("AnalysisRunnerService", "Run",
			"CabgenPipeline")...,
	)

	return nil
}
