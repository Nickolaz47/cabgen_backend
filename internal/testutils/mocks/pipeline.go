package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
)

type MockToolRunner struct {
	RunFunc func(ctx context.Context, args []string) (string, error)

	BuildBlastXCmdFunc        func(blastDB, inputFile, outputFile string) []string
	BuildFastQCCmdFunc        func(fastqcCmd, read1, read2, outputDir string) []string
	BuildUnicyclerCmdFunc     func(unicyclerCmd, read1, read2, outputDir, threads, spadesPath string) []string
	BuildProkkaCmdFunc        func(prokkaCmd, outputDir, prefix, assemblyPath, threads string) []string
	BuildCheckMLineageCmdFunc func(checkmCmd, inputDir, outputDir, threads string) []string
	BuildCheckMQACmdFunc      func(checkmCmd, checkmDir, sample, threads string) []string
	BuildKraken2CmdFunc       func(krakenCmd, dbPath, outputDir, threads, assemblyPath string) []string
	BuildSplitterCmdFunc      func(threads, inputFile, outputFilePrefix string) []string
	BuildFastANICmdFunc       func(fastaniCmd, query, refList, output, threads string) []string
	BuildAbricateCmdFunc      func(abricateCmd, db, inputFile, outputFile, threads string) []string
	BuildMLSTCmdFunc          func(mlstCmd, threads, assemblyPath, outputFile string) []string
}

func (m *MockToolRunner) Run(ctx context.Context, args []string) (string, error) {
	if m.RunFunc != nil {
		return m.RunFunc(ctx, args)
	}
	return "", nil
}

func (m *MockToolRunner) BuildBlastXCmd(blastDB, inputFile,
	outputFile string) []string {
	if m.BuildBlastXCmdFunc != nil {
		return m.BuildBlastXCmdFunc(blastDB, inputFile, outputFile)
	}
	return nil
}

func (m *MockToolRunner) BuildFastQCCmd(fastqcCmd, read1, read2,
	outputDir string) []string {
	if m.BuildFastQCCmdFunc != nil {
		return m.BuildFastQCCmdFunc(fastqcCmd, read1, read2, outputDir)
	}
	return nil
}

func (m *MockToolRunner) BuildUnicyclerCmd(unicyclerCmd, read1, read2,
	outputDir, threads, spadesPath string) []string {
	if m.BuildUnicyclerCmdFunc != nil {
		return m.BuildUnicyclerCmdFunc(unicyclerCmd, read1, read2, outputDir,
			threads, spadesPath)
	}
	return nil
}

func (m *MockToolRunner) BuildProkkaCmd(prokkaCmd, outputDir, prefix,
	assemblyPath, threads string) []string {
	if m.BuildProkkaCmdFunc != nil {
		return m.BuildProkkaCmdFunc(prokkaCmd, outputDir, prefix, assemblyPath,
			threads)
	}
	return nil
}

func (m *MockToolRunner) BuildCheckMLineageCmd(checkmCmd, inputDir,
	outputDir, threads string) []string {
	if m.BuildCheckMLineageCmdFunc != nil {
		return m.BuildCheckMLineageCmdFunc(checkmCmd, inputDir, outputDir,
			threads)
	}
	return nil
}

func (m *MockToolRunner) BuildCheckMQACmd(checkmCmd, checkmDir,
	sample, threads string) []string {
	if m.BuildCheckMQACmdFunc != nil {
		return m.BuildCheckMQACmdFunc(checkmCmd, checkmDir, sample, threads)
	}
	return nil
}

func (m *MockToolRunner) BuildKraken2Cmd(krakenCmd, dbPath, outputDir,
	threads, assemblyPath string) []string {
	if m.BuildKraken2CmdFunc != nil {
		return m.BuildKraken2CmdFunc(krakenCmd, dbPath, outputDir, threads,
			assemblyPath)
	}
	return nil
}

func (m *MockToolRunner) BuildSplitterCmd(threads, inputFile,
	outputFilePrefix string) []string {
	if m.BuildSplitterCmdFunc != nil {
		return m.BuildSplitterCmdFunc(threads, inputFile, outputFilePrefix)
	}
	return nil
}

func (m *MockToolRunner) BuildFastANICmd(fastaniCmd, query, refList,
	output, threads string) []string {
	if m.BuildFastANICmdFunc != nil {
		return m.BuildFastANICmdFunc(fastaniCmd, query, refList, output,
			threads)
	}
	return nil
}

func (m *MockToolRunner) BuildAbricateCmd(abricateCmd, db, inputFile,
	outputFile, threads string) []string {
	if m.BuildAbricateCmdFunc != nil {
		return m.BuildAbricateCmdFunc(abricateCmd, db, inputFile, outputFile,
			threads)
	}
	return nil
}

func (m *MockToolRunner) BuildMLSTCmd(mlstCmd, threads, assemblyPath,
	outputFile string) []string {
	if m.BuildMLSTCmdFunc != nil {
		return m.BuildMLSTCmdFunc(mlstCmd, threads, assemblyPath, outputFile)
	}
	return nil
}

type MockCabgenPipeline struct {
	Config        pipeline.ToolsConfig
	RunFastQCFunc func(ctx context.Context, read1, read2,
		outputDir string) (string, string, error)
	RunUnicyclerFunc func(ctx context.Context, threads int,
		read1, read2, spadesPath, outputDir string) (string, error)
	RunProkkaFunc func(ctx context.Context, threads int,
		assembly, outputDir string) error
	RunCheckMFunc func(ctx context.Context, threads int, sample,
		assemblyDir, outputDir string) (*pipeline.CheckMResult, error)
	RunKraken2Func func(ctx context.Context, threads int,
		assembly, outputDir string) (*pipeline.KrakenSpecies,
		*pipeline.KrakenSpecies, error)
	RunBlastXFunc func(ctx context.Context, query, DB,
		outputFile string) error
	RunAbricateFunc func(ctx context.Context, threads int, db, input,
		outputFile string) error
	ProcessSpeciesFunc func(ctx context.Context, threads int,
		sampleID, mostCommon, assemblyPath, outputDir string) (
		*pipeline.SpeciesResult, error)
}

func (m *MockCabgenPipeline) GetConfig() *pipeline.ToolsConfig {
	return &m.Config
}

func (m *MockCabgenPipeline) RunFastQC(ctx context.Context, read1, read2,
	outputDir string) (string, string, error) {
	if m.RunFastQCFunc != nil {
		return m.RunFastQCFunc(ctx, read1, read2, outputDir)
	}
	return "fastqc1.html", "fastqc2.html", nil
}

func (m *MockCabgenPipeline) RunUnicycler(ctx context.Context, threads int,
	read1, read2, spadesPath, outputDir string) (string, error) {
	if m.RunUnicyclerFunc != nil {
		return m.RunUnicyclerFunc(ctx, threads, read1, read2, spadesPath,
			outputDir)
	}
	return "assembly.fasta", nil
}

func (m *MockCabgenPipeline) RunProkka(ctx context.Context, threads int,
	assembly, outputDir string) error {
	if m.RunProkkaFunc != nil {
		return m.RunProkkaFunc(ctx, threads, assembly, outputDir)
	}
	return nil
}

func (m *MockCabgenPipeline) RunCheckM(ctx context.Context, threads int,
	sample, assemblyDir, outputDir string) (*pipeline.CheckMResult, error) {
	if m.RunCheckMFunc != nil {
		return m.RunCheckMFunc(ctx, threads, sample, assemblyDir, outputDir)
	}
	return &pipeline.CheckMResult{
		Completeness: "99.5", Contamination: "0.5",
		GenomeSize: "5000000", N50: "100000",
	}, nil
}

func (m *MockCabgenPipeline) RunKraken2(ctx context.Context, threads int,
	assembly, outputDir string) (*pipeline.KrakenSpecies,
	*pipeline.KrakenSpecies, error) {
	if m.RunKraken2Func != nil {
		return m.RunKraken2Func(ctx, threads, assembly, outputDir)
	}
	return &pipeline.KrakenSpecies{Name: "Escherichia coli", Count: 100},
		&pipeline.KrakenSpecies{Name: "Klebsiella pneumoniae", Count: 5}, nil
}

func (m *MockCabgenPipeline) RunBlastX(ctx context.Context, query, DB,
	outputFile string) error {
	if m.RunBlastXFunc != nil {
		return m.RunBlastXFunc(ctx, query, DB, outputFile)
	}
	return nil
}

func (m *MockCabgenPipeline) RunAbricate(ctx context.Context, threads int,
	db, input, outputFile string) error {
	if m.RunAbricateFunc != nil {
		return m.RunAbricateFunc(ctx, threads, db, input, outputFile)
	}
	return nil
}

func (m *MockCabgenPipeline) ProcessSpecies(ctx context.Context, threads int,
	sampleID, mostCommon, assemblyPath, outputDir string) (
	*pipeline.SpeciesResult, error) {
	if m.ProcessSpeciesFunc != nil {
		return m.ProcessSpeciesFunc(ctx, threads, sampleID, mostCommon,
			assemblyPath, outputDir)
	}
	return &pipeline.SpeciesResult{
		DisplayName: "Escherichia coli", MLSTSpecies: "ecoli (ST: 131)",
	}, nil
}
