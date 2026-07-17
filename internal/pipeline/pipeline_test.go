package pipeline_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
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
		PoliDbEntero:        "/blast/poli/proteins_Ecloacae_poli.fasta",
		PoliDbAcineto:      "/blast/poli/proteins_acineto_poli.fasta",
		OtherDbPseudo:      "/blast/other/proteins_outrasMut_pseudo.fasta",
		OtherDbKleb:        "/blast/other/proteins_outrasMut_kleb.fasta",
		OtherDbEntero:       "/blast/other/proteins_outrasMut_Ecloacae.fasta",
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

// Duplicated from mutations_test.go / kraken_test.go (package pipeline) since
// this file uses the external test package (pipeline_test) to avoid an import
// cycle with testutils/mocks.
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

func krakenLine(seqID, taxon string) string {
	return "C\t" + seqID + "\t" + taxon + "\t|0:0|\n"
}

func TestNewCabgenPipeline(t *testing.T) {
	p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
		defaultConfig())
	assert.NotNil(t, p)
}

func TestRunFastQC(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
		html1, html2, err := p.RunFastQC(context.Background(),
			"/data/r1.fq", "/data/r2.fq", "/out")
		assert.NoError(t, err)
		assert.Equal(t, "/out/r1.fq_fastqc.html", html1)
		assert.Equal(t, "/out/r2.fq_fastqc.html", html2)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig())
		_, _, err := p.RunFastQC(context.Background(), "r1", "r2", "/out")
		assert.Error(t, err)
	})
}

func TestRunUnicycler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
		path, err := p.RunUnicycler(context.Background(), 4, "r1", "r2",
			"/spades", "/out")
		assert.NoError(t, err)
		assert.Equal(t, "/out/assembly.fasta", path)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig())
		_, err := p.RunUnicycler(context.Background(), 4, "r1", "r2", "",
			"/out")
		assert.Error(t, err)
	})
}

func TestRunProkka(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
		err := p.RunProkka(context.Background(), 8, "contigs.fa", "/out")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig())
		err := p.RunProkka(context.Background(), 8, "contigs.fa", "/out")
		assert.Error(t, err)
	})
}

func TestRunBlastX(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
		err := p.RunBlastX(context.Background(), "contigs.fa", "/db",
			"out.txt")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: errorRun},
			defaultConfig())
		err := p.RunBlastX(context.Background(), "contigs.fa", "/db",
			"out.txt")
		assert.Error(t, err)
	})
}

func TestRunCheckM(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		outDir := t.TempDir()
		writeFile(t, filepath.Join(outDir, "s1_results"),
			"Bin Id\tML\tG\tM\tMS\tComp\tCont\tC\tGS\tN\tN100\tN50\n"+
				"s1\tF\t5\t10\t5\t98.5\t0.5\t3\t3500000\t0\t0\t25000\n")

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
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
			defaultConfig())
		result, err := p.RunCheckM(context.Background(), 4, "s1", "/in",
			"/out")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - QA Fails", func(t *testing.T) {
		calls := 0
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{
			RunFunc: func(_ context.Context, _ []string) (string, error) {
				calls++
				if calls == 2 {
					return "", fmt.Errorf("qa failed")
				}
				return "", nil
			},
		}, defaultConfig())

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
		writeFile(t, filepath.Join(outDir, "out_kraken"),
			krakenLine("r1", "Escherichia coli")+
				krakenLine("r2", "Klebsiella pneumoniae"))

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
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
			defaultConfig())
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
			defaultConfig())
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
			defaultConfig())
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
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
			organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
			organismMockContent)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Klebsiella pneumoniae", "contigs.fa", outDir)
		assert.NoError(t, err)
		assert.Contains(t, result.OtherMutations, "GyrA:A3C,")
	})

	t.Run("Success - Single Word Species Name", func(t *testing.T) {
		outDir := t.TempDir()
		sampleID := "s1"
		writeFile(t, filepath.Join(outDir, sampleID+"_blastPoli"),
			organismMockContent)
		writeFile(t, filepath.Join(outDir, sampleID+"_blastOther"),
			organismMockContent)

		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{RunFunc: successRun},
			defaultConfig())
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
			RunFunc: func(_ context.Context, _ []string) (string, error) {
				return mlstPath, nil
			},
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
		p := pipeline.NewCabgenPipeline(&mocks.MockToolRunner{
			RunFunc: func(_ context.Context, _ []string) (string, error) {
				calls++
				if calls >= 3 {
					return "", fmt.Errorf("blastx failed")
				}
				return "", nil
			},
		}, defaultConfig())

		// Enterobacter cloacae triggers the Enterobacter branch, which
		// runs MLST, FastANI, then BlastX (poli first). The mock fails on
		// call 3, which is the BlastX poli invocation.
		result, err := p.ProcessSpecies(context.Background(), 4, sampleID,
			"Enterobacter cloacae", "contigs.fa", outDir)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "blastx failed")
	})
}