package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPipelineRunner struct {
	runFunc func(ctx context.Context, args []string) (string, error)
}

func (m *mockPipelineRunner) Run(ctx context.Context, args []string) (string, error) {
	return m.runFunc(ctx, args)
}
func (m *mockPipelineRunner) BuildBlastXCmd(_, _, _ string) []string             { return nil }
func (m *mockPipelineRunner) BuildFastQCCmd(_, _, _, _ string) []string          { return nil }
func (m *mockPipelineRunner) BuildUnicyclerCmd(_, _, _, _, _, _ string) []string   { return nil }
func (m *mockPipelineRunner) BuildProkkaCmd(_, _, _, _, _ string) []string        { return nil }
func (m *mockPipelineRunner) BuildCheckMLineageCmd(_, _, _, _ string) []string    { return nil }
func (m *mockPipelineRunner) BuildCheckMQACmd(_, _, _, _ string) []string         { return nil }
func (m *mockPipelineRunner) BuildKraken2Cmd(_, _, _, _, _ string) []string        { return nil }
func (m *mockPipelineRunner) BuildSplitterCmd(_, _, _ string) []string             { return nil }
func (m *mockPipelineRunner) BuildFastANICmd(_, _, _, _, _ string) []string        { return nil }
func (m *mockPipelineRunner) BuildAbricateCmd(_, _, _, _, _ string) []string       { return nil }
func (m *mockPipelineRunner) BuildMLSTCmd(_, _, _, _ string) []string              { return nil }

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
}

func defaultConfig() ToolsConfig {
	return ToolsConfig{
		FastQCPath:       "fastqc",
		UnicyclerPath:    "unicycler",
		SpadesPath:       "/spades",
		ProkkaPath:       "prokka",
		CheckMPath:       "checkm",
		Kraken2Path:      "kraken2",
		KrakenDBPath:     "/db",
		FastANIPath:      "fastani",
		FastANIRefsPath:  "/refs",
		AbricatePath:     "abricate",
		MLSTPath:         "mlst",
		BlastPoliDBPath:  "/poliDB",
		BlastOtherDBPath: "/otherDB",
	}
}

var successRun = funcRun("")
var errorRun = funcRunErr(fmt.Errorf("command failed"))

func funcRun(stdout string) func(context.Context, []string) (string, error) {
	return func(_ context.Context, _ []string) (string, error) { return stdout, nil }
}

func funcRunErr(err error) func(context.Context, []string) (string, error) {
	return func(_ context.Context, _ []string) (string, error) { return "", err }
}

func TestNewCabgenPipeline(t *testing.T) {
	p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
	assert.NotNil(t, p)
}

func TestRunFastQC(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		html1, html2, err := p.RunFastQC(context.Background(), "/data/r1.fq", "/data/r2.fq", "/out")
		assert.NoError(t, err)
		assert.Equal(t, "/out/r1.fq_fastqc.html", html1)
		assert.Equal(t, "/out/r2.fq_fastqc.html", html2)
	})

	t.Run("Error", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: errorRun}, defaultConfig())
		_, _, err := p.RunFastQC(context.Background(), "r1", "r2", "/out")
		assert.Error(t, err)
	})
}

func TestRunUnicycler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		path, err := p.RunUnicycler(context.Background(), 4, "r1", "r2", "/spades", "/out")
		assert.NoError(t, err)
		assert.Equal(t, "/out/assembly.fasta", path)
	})

	t.Run("Error", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: errorRun}, defaultConfig())
		_, err := p.RunUnicycler(context.Background(), 4, "r1", "r2", "", "/out")
		assert.Error(t, err)
	})
}

func TestRunProkka(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		err := p.RunProkka(context.Background(), 8, "contigs.fa", "/out")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: errorRun}, defaultConfig())
		err := p.RunProkka(context.Background(), 8, "contigs.fa", "/out")
		assert.Error(t, err)
	})
}

func TestRunBlastX(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		err := p.RunBlastX(context.Background(), "contigs.fa", "/db", "out.txt")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: errorRun}, defaultConfig())
		err := p.RunBlastX(context.Background(), "contigs.fa", "/db", "out.txt")
		assert.Error(t, err)
	})
}

func TestRunCheckM(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		outDir := t.TempDir()
		writeFile(t, filepath.Join(outDir, "s1_results"),
			"Bin Id\tML\tG\tM\tMS\tComp\tCont\tC\tGS\tN\tN100\tN50\n"+
				"s1\tF\t5\t10\t5\t98.5\t0.5\t3\t3500000\t0\t0\t25000\n")

		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		result, err := p.RunCheckM(context.Background(), 4, "s1", "/in", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "98.5", result.Completeness)
		assert.Equal(t, "0.5", result.Contamination)
		assert.Equal(t, "3500000", result.GenomeSize)
		assert.Equal(t, "25000", result.N50)
	})

	t.Run("Error - Lineage Fails", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: errorRun}, defaultConfig())
		result, err := p.RunCheckM(context.Background(), 4, "s1", "/in", "/out")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - QA Fails", func(t *testing.T) {
		calls := 0
		p := NewCabgenPipeline(&mockPipelineRunner{
			runFunc: func(_ context.Context, _ []string) (string, error) {
				calls++
				if calls == 2 {
					return "", fmt.Errorf("qa failed")
				}
				return "", nil
			},
		}, defaultConfig())

		result, err := p.RunCheckM(context.Background(), 4, "s1", "/in", "/out")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "qa failed")
	})
}

func TestRunKraken2(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		outDir := t.TempDir()
		writeFile(t, filepath.Join(outDir, "out_kraken"),
			krakenLine("r1", "Escherichia coli")+
				krakenLine("r2", "Klebsiella pneumoniae"))

		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		first, second, err := p.RunKraken2(context.Background(), 4, "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.NotNil(t, second)
		assert.Equal(t, "Klebsiella pneumoniae", second.Name)
	})

	t.Run("Error", func(t *testing.T) {
		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: errorRun}, defaultConfig())
		first, second, err := p.RunKraken2(context.Background(), 4, "contigs.fa", "/out")
		assert.Error(t, err)
		assert.Nil(t, first)
		assert.Nil(t, second)
	})
}

func TestProcessSpecies(t *testing.T) {
	t.Run("Success - Non-matched Species", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"), organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"), organismMockContent)

		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Staphylococcus aureus", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Staphylococcus aureus", result.DisplayName)
		assert.Empty(t, result.MLSTSpecies)
		assert.Empty(t, result.PoliMutations)
		assert.Empty(t, result.OtherMutations)
	})

	t.Run("Success - Acinetobacter Finds Mutations", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"), organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"), organismMockContent)

		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Acinetobacter baumannii", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Acinetobacter baumannii", result.DisplayName)
		assert.Contains(t, result.OtherMutations, "GyrA:A3C,")
		assert.Contains(t, result.PoliMutations, "PmrA:A3C,")
	})

	t.Run("Success - Klebsiella Finds Mutations", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"), organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"), organismMockContent)

		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Klebsiella pneumoniae", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.Contains(t, result.OtherMutations, "GyrA:A3C,")
	})

	t.Run("Success - Single Word Species Name", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"), organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"), organismMockContent)

		p := NewCabgenPipeline(&mockPipelineRunner{runFunc: successRun}, defaultConfig())
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Acinetobacter", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.Equal(t, "Acinetobacter", result.DisplayName)
	})

	t.Run("Success - MLST Parsed", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		mlstPath := filepath.Join(outDir, "mlst.csv")
		writeFile(t, mlstPath, "contigs.fa,abaumannii,ST2,oxa0001,ompA0001\n")
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"), organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"), organismMockContent)

		p := NewCabgenPipeline(&mockPipelineRunner{
			runFunc: funcRun(mlstPath),
		}, defaultConfig())
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Acinetobacter baumannii", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.Equal(t, "abaumannii (ST: ST2)", result.MLSTSpecies)
	})

	t.Run("Error - BlastX Poli Fails", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		calls := 0
		p := NewCabgenPipeline(&mockPipelineRunner{
			runFunc: func(_ context.Context, _ []string) (string, error) {
				calls++
				if calls >= 3 {
					return "", fmt.Errorf("blastx failed")
				}
				return "", nil
			},
		}, defaultConfig())

		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Escherichia coli", "contigs.fa", outDir)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "blastx failed")
	})
}