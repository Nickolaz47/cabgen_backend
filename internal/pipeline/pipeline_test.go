package pipeline_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
}

func defaultConfig() pipeline.ToolsConfig {
	return pipeline.ToolsConfig{
		FastQCPath:         "fastqc",
		UnicyclerPath:      "unicycler",
		SpadesPath:         "/spades",
		CheckMPath:         "checkm",
		Kraken2Path:        "kraken2",
		KrakenDBPath:       "/db",
		FastANIPath:        "fastani",
		AbricatePath:       "abricate",
		MLSTPath:           "mlst",
		ResfinderDBPath:    "/resfinder_db",
		PoliDbPseudo:       "/blast/poli/proteins_pseudo_poli.fasta",
		PoliDbKleb:         "/blast/poli/proteins_kleb_poli.fasta",
		PoliDbEntero:       "/blast/poli/proteins_Ecloacae_poli.fasta",
		PoliDbAcineto:      "/blast/poli/proteins_acineto_poli.fasta",
		OtherDbPseudo:      "/blast/other/proteins_outrasMut_pseudo.fasta",
		OtherDbKleb:        "/blast/other/proteins_outrasMut_kleb.fasta",
		OtherDbEntero:      "/blast/other/proteins_outrasMut_Ecloacae.fasta",
		OtherDbAcineto:     "/blast/other/proteins_outrasMut_acineto.fasta",
		FastaniListKleb:    "/fastani/kleb_database/lista-kleb",
		FastaniListEntero:  "/fastani/fastANI/list_entero",
		FastaniListAcineto: "/fastani/fastANI_acineto/list-acineto",
	}
}

func successRun(_ context.Context, _ []string) (string, error) {
	return "", nil
}

func errorRun(_ context.Context, _ []string) (string, error) {
	return "", fmt.Errorf("command failed")
}

const organismMockContent = `
> GyrA|
Length=100
Identities = 95/95 (95%)
Query  1   ATCG 4
           || |
Sbjct  1   ATAG 4

> PmrA|
Length=100
Identities = 95/95 (95%)
Query  1   ATCG 4
           || |
Sbjct  1   ATAG 4
`

func krakenReportLine(name string, cladeReads int) string {
	return fmt.Sprintf("1.00\t%d\t%d\tS\t0\t%s\n", cladeReads, cladeReads,
		name)
}

func TestNewCabgenPipeline(t *testing.T) {
	p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
		defaultConfig(), nil)
	assert.NotNil(t, p)
}

func TestRunFastQC(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		html1, html2, err := p.RunFastQC(context.Background(),
			"/data/r1.fq", "/data/r2.fq", "/out")
		assert.NoError(t, err)
		assert.Equal(t, "/out/r1_fastqc.html", html1)
		assert.Equal(t, "/out/r2_fastqc.html", html2)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig(), nil)
		_, _, err := p.RunFastQC(context.Background(), "r1", "r2", "/out")
		assert.Error(t, err)
	})
}

func TestRunUnicycler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		outDir := t.TempDir()
		assemblyFile := filepath.Join(outDir, "assembly.fasta")
		os.WriteFile(assemblyFile, []byte(">seq1\nATCG\n"), 0644)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		path, err := p.RunUnicycler(context.Background(), 4, "r1", "r2",
			"/spades", outDir, "A01_assembly.fasta")
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(outDir, "A01_assembly.fasta"), path)
		assert.FileExists(t, path)
		assert.NoFileExists(t, assemblyFile)
	})

	t.Run("Success - File Missing Falls Back", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		path, err := p.RunUnicycler(context.Background(), 4, "r1", "r2",
			"/spades", "/out", "A01_assembly.fasta")
		assert.NoError(t, err)
		assert.Equal(t, "/out/assembly.fasta", path)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig(), nil)
		_, err := p.RunUnicycler(context.Background(), 4, "r1", "r2", "",
			"/out", "A01_assembly.fasta")
		assert.Error(t, err)
	})
}

func TestRunProkka(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		err := p.RunProkka(context.Background(), 8, "contigs.fa", "/out")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig(), nil)
		err := p.RunProkka(context.Background(), 8, "contigs.fa", "/out")
		assert.Error(t, err)
	})
}

func TestRunBlastX(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		err := p.RunBlastX(context.Background(), "contigs.fa", "/db",
			"out.txt")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig(), nil)
		err := p.RunBlastX(context.Background(), "contigs.fa", "/db",
			"out.txt")
		assert.Error(t, err)
	})
}

func TestRunCheckM(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		outDir := t.TempDir()
		writeFile(t, filepath.Join(outDir, "s1_results"),
			"Bin Id\tML\tG\tM\tMS\tComp\tCont\tSH\tGS\tGC\tC\tS\tN\tN50\n"+
				"s1\tF\t5\t10\t5\t98.5\t0.5\t0\t3500000\t37.5\t3\t2\t0\t25000\n")

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		result, err := p.RunCheckM(context.Background(), 4, "s1", "/in",
			outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "98.5", result.Completeness)
		assert.Equal(t, "0.5", result.Contamination)
		assert.Equal(t, "3500000", result.GenomeSize)
		assert.Equal(t, "25000", result.N50)
	})

	t.Run("Error - Lineage Fails", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig(), nil)
		result, err := p.RunCheckM(context.Background(), 4, "s1", "/in",
			"/out")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - QA Fails", func(t *testing.T) {
		calls := 0
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{
			RunFunc: func(ctx context.Context, args []string) (string, error) {
				calls++
				if calls == 2 {
					return "", fmt.Errorf("qa failed")
				}
				return "", nil
			},
		}, defaultConfig(), nil)

		result, err := p.RunCheckM(context.Background(), 4, "s1", "/in",
			"/out")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "qa failed")
	})
}

func TestRunKraken2(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		outDir := t.TempDir()
		writeFile(t, filepath.Join(outDir, "report_kraken"),
			krakenReportLine("Escherichia coli", 1)+
				krakenReportLine("Klebsiella pneumoniae", 1))

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		first, second, err := p.RunKraken2(context.Background(), 4,
			"contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.NotNil(t, second)
		assert.Equal(t, "Klebsiella pneumoniae", second.Name)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig(), nil)
		first, second, err := p.RunKraken2(context.Background(), 4,
			"contigs.fa", "/out")
		assert.Error(t, err)
		assert.Nil(t, first)
		assert.Nil(t, second)
	})
}

func TestProcessSpecies(t *testing.T) {
	t.Run("Success - Non-matched Species", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
			organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
			organismMockContent)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
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
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
			organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
			organismMockContent)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Acinetobacter baumannii", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Acinetobacter baumannii", result.DisplayName)
		assert.Contains(t, result.OtherMutations, "GyrA:A3C")
		assert.Contains(t, result.PoliMutations, "PmrA:A3C")
	})

	t.Run("Success - Klebsiella Finds Mutations", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
			organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
			organismMockContent)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Klebsiella pneumoniae", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.Contains(t, result.OtherMutations, "GyrA:A3C")
	})

	t.Run("Success - Single Word Species Name", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
			organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
			organismMockContent)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), nil)
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Acinetobacter", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.Equal(t, "Acinetobacter", result.DisplayName)
	})

	t.Run("Success - MLST Parsed", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		mlstPath := filepath.Join(outDir, "mlst.csv")
		writeFile(t, mlstPath,
			"contigs.fa,abaumannii,ST2,oxa0001,ompA0001\n")
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
			organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
			organismMockContent)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{
			RunFunc: func(ctx context.Context, args []string) (string, error) {
				return mlstPath, nil
			},
		}, defaultConfig(), nil)
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Acinetobacter baumannii", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.Equal(t, "abaumannii (ST: ST2)", result.MLSTSpecies)
	})

	t.Run("Success - MLST Skips When Scheme And ST Are Dash",
		func(t *testing.T) {
			outDir := t.TempDir()
			sampleID := "s1"
			mlstPath := filepath.Join(outDir, "mlst.csv")
			writeFile(t, mlstPath,
				"contigs.fa,-,-,oxa0001,ompA0001\n")
			writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
				organismMockContent)
			writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
				organismMockContent)

			p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{
				RunFunc: func(ctx context.Context, args []string) (
					string, error) {
					return mlstPath, nil
				},
			}, defaultConfig(), nil)
			result, err := p.ProcessSpecies(context.Background(), 4,
				sampleID, "Acinetobacter baumannii", "contigs.fa",
				outDir)
			assert.NoError(t, err)
			assert.Empty(t, result.MLSTSpecies)
		})

	t.Run("Success - BlastX Poli Fails Returns Partial Result", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		calls := 0
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{
			RunFunc: func(ctx context.Context, args []string) (string, error) {
				calls++
				if calls >= 3 {
					return "", fmt.Errorf("blastx failed")
				}
				return "", nil
			},
		}, defaultConfig(), nil)

		// Enterobacter cloacae triggers the Enterobacter branch, which
		// runs MLST, FastANI, then BlastX (poli first). The mock fails on
		// call 3, which is the BlastX poli invocation.
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Enterobacter cloacae", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Enterobacter cloacae", result.DisplayName)
		assert.Empty(t, result.PoliMutations)
		assert.Empty(t, result.OtherMutations)
	})
}

func TestProcessSpeciesLogging(t *testing.T) {
	t.Run("Logs debug when no species matches", func(t *testing.T) {
		outDir := t.TempDir()
		logger, logs := testutils.NewMockLogger(zapcore.DebugLevel)

		p := pipeline.NewCabgenPipeline(
			&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), logger)

		result, err := p.ProcessSpecies(context.Background(), 4, "s1",
			"Staphylococcus aureus", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Staphylococcus aureus", result.DisplayName)

		found := false
		for _, entry := range logs.All() {
			if entry.Level == zapcore.DebugLevel &&
				entry.Message == "Species did not match any known genus, skipping BlastX/FastANI" {
				found = true
			}
		}
		assert.True(t, found, "expected debug log about unknown genus")
	})

	t.Run("Logs warning when FastANI ref list is empty", func(t *testing.T) {
		outDir := t.TempDir()
		logger, logs := testutils.NewMockLogger(zapcore.WarnLevel)

		cfg := defaultConfig()
		cfg.FastaniListAcineto = ""

		p := pipeline.NewCabgenPipeline(
			&mocks.MockToolRunner{RunFunc: successRun},
			cfg, logger)

		result, err := p.ProcessSpecies(context.Background(), 4, "s1",
			"Acinetobacter baumannii", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		found := false
		for _, entry := range logs.All() {
			if entry.Level == zapcore.WarnLevel &&
				entry.Message == "Matched genus but FASTANI ref list not configured, skipping FastANI" {
				found = true
			}
		}
		assert.True(t, found, "expected warning log about missing FastANI ref list")
	})

	t.Run("Logs error when FastANI execution fails", func(t *testing.T) {
		outDir := t.TempDir()
		logger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		calls := 0
		p := pipeline.NewCabgenPipeline(
			&mocks.MockToolRunner{
				RunFunc: func(ctx context.Context, args []string) (string, error) {
					calls++
					if calls == 2 {
						return "", fmt.Errorf("fastani binary not found")
					}
					return "", nil
				},
			}, defaultConfig(), logger)

		result, err := p.ProcessSpecies(context.Background(), 4, "s1",
			"Acinetobacter baumannii", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Acinetobacter baumannii", result.DisplayName)

		found := false
		for _, entry := range logs.All() {
			if entry.Level == zapcore.ErrorLevel &&
				entry.Message == "FastANI failed" {
				found = true
			}
		}
		assert.True(t, found, "expected error log about FastANI failure")
	})

	t.Run("Logs warning when ParseFastANI fails", func(t *testing.T) {
		outDir := t.TempDir()
		logger, logs := testutils.NewMockLogger(zapcore.WarnLevel)

		p := pipeline.NewCabgenPipeline(
			&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig(), logger)

		result, err := p.ProcessSpecies(context.Background(), 4, "s1",
			"Acinetobacter baumannii", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Acinetobacter baumannii", result.DisplayName)

		found := false
		for _, entry := range logs.All() {
			if entry.Level == zapcore.WarnLevel &&
				entry.Message == "FastANI output parse failed" {
				found = true
			}
		}
		assert.True(t, found, "expected warning log about FastANI parse failure")
	})
}
