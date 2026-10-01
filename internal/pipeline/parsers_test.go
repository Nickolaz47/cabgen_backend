package pipeline

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockParserFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "parser_mock_*.txt")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

func TestParseCheckM(t *testing.T) {
	const header = "Bin Id\tMarker lineage\t# genomes\t# markers\t# marker sets\tCompleteness\tContamination\tStrain heterogeneity\tGenome size (bp)\t# ambiguous bases\t# scaffolds\t# contigs\tN50 (scaffolds)\tN50 (contigs)\n"

	t.Run("Success - Valid CheckM Output", func(t *testing.T) {
		content := header +
			"cabgen4_assembly\tf__Moraxellaceae (UID4680)\t86\t689\t365\t97.95\t0.66\t50.00\t3651976\t0\t745\t745\t7060\t7060\t4901\t4901\t35268\t35268\t40.0\t2.14\t89.54\t11\t3725\t9\t676\t4\t0\t0\t0\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "97.95", result.Completeness)
		assert.Equal(t, "0.66", result.Contamination)
		assert.Equal(t, "3651976", result.GenomeSize)
		assert.Equal(t, "745", result.Contigs)
		assert.Equal(t, "7060", result.N50)
	})

	t.Run("Success - Header Skipped", func(t *testing.T) {
		content := header +
			"sample1\tFirmicutes\t543\t124\t58\t99.20\t1.05\t0\t4200000\t0\t520\t520\t31000\t31000\t6000\t6000\t45000\t45000\t40.0\t2.0\t90.0\t11\t4000\t5\t300\t2\t0\t0\t0\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "99.20", result.Completeness)
		assert.Equal(t, "1.05", result.Contamination)
		assert.Equal(t, "4200000", result.GenomeSize)
		assert.Equal(t, "520", result.Contigs)
		assert.Equal(t, "31000", result.N50)
	})

	t.Run("Error - Empty File", func(t *testing.T) {
		path := createMockParserFile(t, "")

		result, err := ParseCheckM(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Empty checkm result")
	})

	t.Run("Error - Only Header No Data", func(t *testing.T) {
		path := createMockParserFile(t, header)

		result, err := ParseCheckM(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "No valid data found in checkm result")
	})

	t.Run("Error - Data Line With Fewer Than 14 Fields", func(t *testing.T) {
		content := header +
			"sample1\tFirmicutes\t543\t124\t58\t97.95\t0.66\t50.00\t3651976\t0\t745\t745\t7060\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "No valid data found in checkm result")
	})

	t.Run("Success - Blank Lines Skipped", func(t *testing.T) {
		content := header +
			"\n" +
			"\n" +
			"cabgen4_assembly\tf__Moraxellaceae (UID4680)\t86\t689\t365\t97.95\t0.66\t50.00\t3651976\t0\t745\t745\t7060\t7060\t4901\t4901\t35268\t35268\t40.0\t2.14\t89.54\t11\t3725\t9\t676\t4\t0\t0\t0\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "97.95", result.Completeness)
		assert.Equal(t, "745", result.Contigs)
		assert.Equal(t, "7060", result.N50)
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		result, err := ParseCheckM("nonexistent.txt")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Failed to open checkm result")
	})
}

func TestParseFastANI(t *testing.T) {
	t.Run("Success - Valid FastANI Output", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/Ecoli_K12.fasta\t99.87\t1500/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.NoError(t, err)
		assert.Equal(t, "Ecoli_K12", result)
	})

	t.Run("Success - Multiple Lines Returns Best ANI", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/Salmonella.fasta\t95.20\t1200/1500\n" +
			"/data/contigs.fa\t/data/ref/Ecoli_K12.fasta\t99.87\t1500/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.NoError(t, err)
		assert.Equal(t, "Ecoli_K12", result)
	})

	t.Run("Success - Ref Without Path", func(t *testing.T) {
		content := "/data/contigs.fa\tKpneumo.fna\t97.50\t1200/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.NoError(t, err)
		assert.Equal(t, "Kpneumo", result)
	})

	t.Run("Success - Ref With Multiple Dots In Filename", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/sample.genome.v2.fasta\t99.87\t1500/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.NoError(t, err)
		assert.Equal(t, "sample", result)
	})

	t.Run("Success - Blank Lines Skipped", func(t *testing.T) {
		content := "\n" +
			"\n" +
			"/data/contigs.fa\t/data/ref/Ecoli_K12.fasta\t99.87\t1500/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.NoError(t, err)
		assert.Equal(t, "Ecoli_K12", result)
	})

	t.Run("Success - Best ANI Below 95% Returns Error", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/Other.fasta\t94.00\t1500/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(),
			"best ANI match below species threshold")
	})

	t.Run("Success - Invalid ANI Value Ignored", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/Bad.fasta\tnot-a-number\t1500/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(),
			"no valid data found in fastani result")
	})

	t.Run("Error - Empty File", func(t *testing.T) {
		path := createMockParserFile(t, "")

		result, err := ParseFastANI(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "no valid data found in fastani result")
	})

	t.Run("Error - Line With Fewer Than 3 Fields", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/Ecoli_K12.fasta\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "no valid data found in fastani result")
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		result, err := ParseFastANI("nonexistent.txt")
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "Failed to open fastani result")
	})
}

func TestParseMLST(t *testing.T) {
	t.Run("Success - Valid MLST CSV Output", func(t *testing.T) {
		content := "contigs.fa,ecoli,131,adek0001,fyhn0001,gyrA0001,icd0001,mdh0001,purA0001,recA0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.Equal(t, "ecoli (ST131)", result)
	})

	t.Run("Success - Quoted Fields In CSV", func(t *testing.T) {
		content := "contigs.fa,abaumannii,2,oxa0001,ompA0001,csuE0001,fkpA0001,rplB0001,gltA0001,gyrB0001,gdhB0001,recA0001,gpi0001,rpoB0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.Equal(t, "abaumannii (ST2)", result)
	})

	t.Run("Success - Multiple Lines Returns First", func(t *testing.T) {
		content := "contigs1.fa,ecoli,131,adek0001,fyhn0001,gyrA0001\n" +
			"contigs2.fa,kpneumo,258,tonB0001,infB0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.Equal(t, "ecoli (ST131)", result)
	})

	t.Run("Success - Blank Lines Skipped", func(t *testing.T) {
		content := "\n" +
			"\n" +
			"contigs.fa,ecoli,131,adek0001,fyhn0001,gyrA0001,icd0001,mdh0001,purA0001,recA0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.Equal(t, "ecoli (ST131)", result)
	})

	t.Run("Success - New ST", func(t *testing.T) {
		content := "contigs.fa,ecoli,-,adek0001,fyhn0001,gyrA0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.Equal(t, "ecoli (New ST)", result)
	})

	t.Run("Success - Not available for this species", func(t *testing.T) {
		content := "contigs.fa,-,-,adek0001,fyhn0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.Equal(t, "Not available for this species", result)
	})

	t.Run("Error - Empty File", func(t *testing.T) {
		path := createMockParserFile(t, "")

		result, err := ParseMLST(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "No valid data found in mlst result")
	})

	t.Run("Error - Line With Fewer Than 3 Fields", func(t *testing.T) {
		content := "contigs.fa,ecoli\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "No valid data found in mlst result")
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		result, err := ParseMLST("nonexistent.txt")
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "Failed to open mlst result")
	})
}

func TestParseFastANI_ScannerError(t *testing.T) {
	t.Run("Error - Line Exceeds Scanner Buffer", func(t *testing.T) {
		hugeLine := "path\t/data/ref/Ecoli.fasta\t99.9\t" + strings.Repeat("a", 128*1024)
		path := createMockParserFile(t, hugeLine)

		result, err := ParseFastANI(path)

		assert.Empty(t, result)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Error reading fastani result")
		assert.ErrorContains(t, err, "token too long")
	})
}

func TestParseMLST_ScannerError(t *testing.T) {
	t.Run("Error - Line Exceeds Scanner Buffer", func(t *testing.T) {
		hugeLine := "contigs.fa,ecoli,131,adek0001\t" + strings.Repeat("a", 128*1024)
		path := createMockParserFile(t, hugeLine)

		result, err := ParseMLST(path)

		assert.Empty(t, result)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Error reading mlst result")
		assert.ErrorContains(t, err, "token too long")
	})
}
