package pipeline

import (
	"os"
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
	// checkm qa -o 2 --tab_table output format (tab-separated):
	// Bin Id  Marker lineage  # genomes  # markers  # marker sets
	// Completeness  Contamination  # contigs  Genome size  # N's
	// # N's per 100 kbp  Scaffold N50
	t.Run("Success - Valid CheckM Output", func(t *testing.T) {
		content := "Bin Id\tMarker lineage\t# genomes\t# markers\t# marker sets\tCompleteness\tContamination\t# contigs\tGenome size\t# N's\t# N's per 100 kbp\tScaffold N50\n" +
			"sample1\tFirmicutes\t543\t124\t58\t98.54\t0.52\t15\t3500000\t0\t0\t25000\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "98.54", result.Completeness)
		assert.Equal(t, "0.52", result.Contamination)
		assert.Equal(t, "3500000", result.GenomeSize)
		assert.Equal(t, "25000", result.N50)
	})

	t.Run("Success - Header Skipped", func(t *testing.T) {
		content := "Bin Id\tMarker lineage\t# genomes\t# markers\t# marker sets\tCompleteness\tContamination\t# contigs\tGenome size\t# N's\t# N's per 100 kbp\tScaffold N50\n" +
			"sample1\tFirmicutes\t543\t124\t58\t99.20\t1.05\t10\t4200000\t0\t0\t30000\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "99.20", result.Completeness)
		assert.Equal(t, "1.05", result.Contamination)
		assert.Equal(t, "4200000", result.GenomeSize)
		assert.Equal(t, "30000", result.N50)
	})

	t.Run("Error - Empty File", func(t *testing.T) {
		path := createMockParserFile(t, "")

		result, err := ParseCheckM(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Empty checkm result")
	})

	t.Run("Error - Only Header No Data", func(t *testing.T) {
		content := "Bin Id\tMarker lineage\t# genomes\t# markers\t# marker sets\tCompleteness\tContamination\t# contigs\tGenome size\t# N's\t# N's per 100 kbp\tScaffold N50\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "No valid data found in checkm result")
	})

	t.Run("Error - Data Line With Fewer Than 12 Fields", func(t *testing.T) {
		content := "Bin Id\tMarker lineage\t# genomes\t# markers\t# marker sets\tCompleteness\tContamination\n" +
			"sample1\tFirmicutes\t543\t124\t58\t98.54\t0.52\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "No valid data found in checkm result")
	})

	t.Run("Error - Blank Lines Skipped", func(t *testing.T) {
		content := "Bin Id\tMarker lineage\t# genomes\t# markers\t# marker sets\tCompleteness\tContamination\t# contigs\tGenome size\t# N's\t# N's per 100 kbp\tScaffold N50\n" +
			"\n" +
			"\n" +
			"sample1\tFirmicutes\t543\t124\t58\t98.54\t0.52\t15\t3500000\t0\t0\t25000\n"
		path := createMockParserFile(t, content)

		result, err := ParseCheckM(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "98.54", result.Completeness)
		assert.Equal(t, "25000", result.N50)
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		result, err := ParseCheckM("nonexistent.txt")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Failed to open checkm result")
	})
}

func TestParseFastANI(t *testing.T) {
	// fastani -q <query> --rl <refList> -o <output> output format (tab-separated):
	// query  ref  ANI  fragments_matched/total_fragments
	// Parser extracts filename from ref path, strips extension
	t.Run("Success - Valid FastANI Output", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/Ecoli_K12.fasta\t99.87\t1500/1500\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.NoError(t, err)
		assert.Equal(t, "Ecoli_K12", result)
	})

	t.Run("Success - Multiple Lines Returns First", func(t *testing.T) {
		content := "/data/contigs.fa\t/data/ref/Ecoli_K12.fasta\t99.87\t1500/1500\n" +
			"/data/contigs.fa\t/data/ref/Salmonella.fasta\t95.20\t1200/1500\n"
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

	t.Run("Error - Empty File", func(t *testing.T) {
		path := createMockParserFile(t, "")

		result, err := ParseFastANI(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "No valid data found in fastani result")
	})

	t.Run("Error - Line With Fewer Than 2 Fields", func(t *testing.T) {
		content := "/data/contigs.fa\n"
		path := createMockParserFile(t, content)

		result, err := ParseFastANI(path)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "No valid data found in fastani result")
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		result, err := ParseFastANI("nonexistent.txt")
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "Failed to open fastani result")
	})
}

func TestParseMLST(t *testing.T) {
	// mlst --csv output format (comma-separated):
	// file,scheme,ST,gene1,gene2,...
	// Parser extracts fields[1]=Scheme, fields[2]=ST
	t.Run("Success - Valid MLST CSV Output", func(t *testing.T) {
		content := "contigs.fa,ecoli,ST131,adek0001,fyhn0001,gyrA0001,icd0001,mdh0001,purA0001,recA0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ecoli", result.Scheme)
		assert.Equal(t, "ST131", result.ST)
	})

	t.Run("Success - Quoted Fields In CSV", func(t *testing.T) {
		content := "\"contigs.fa\",\"abaumannii\",\"ST2\",\"oxa0001\",\"ompA0001\",\"csuE0001\",\"fkpA0001\",\"rplB0001\",\"gltA0001\",\"gyrB0001\",\"gdhB0001\",\"recA0001\",\"gpi0001\",\"rpoB0001\"\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "\"abaumannii\"", result.Scheme)
		assert.Equal(t, "\"ST2\"", result.ST)
	})

	t.Run("Success - Multiple Lines Returns First", func(t *testing.T) {
		content := "contigs1.fa,ecoli,ST131,adek0001,fyhn0001,gyrA0001\n" +
			"contigs2.fa,kpneumo,ST258,tonB0001,infB0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ecoli", result.Scheme)
		assert.Equal(t, "ST131", result.ST)
	})

	t.Run("Success - Blank Lines Skipped", func(t *testing.T) {
		content := "\n" +
			"\n" +
			"contigs.fa,ecoli,ST131,adek0001,fyhn0001,gyrA0001,icd0001,mdh0001,purA0001,recA0001\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ecoli", result.Scheme)
		assert.Equal(t, "ST131", result.ST)
	})

	t.Run("Error - Empty File", func(t *testing.T) {
		path := createMockParserFile(t, "")

		result, err := ParseMLST(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "No valid data found in mlst result")
	})

	t.Run("Error - Line With Fewer Than 3 Fields", func(t *testing.T) {
		content := "contigs.fa,ecoli\n"
		path := createMockParserFile(t, content)

		result, err := ParseMLST(path)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "No valid data found in mlst result")
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		result, err := ParseMLST("nonexistent.txt")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Failed to open mlst result")
	})
}
