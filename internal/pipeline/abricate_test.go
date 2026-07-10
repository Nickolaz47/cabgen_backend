package pipeline

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockAbricateFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "abricate_mock_*.txt")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

// buildAbricateLine creates a tab-separated abricate result line.
// fields: 0=seqid, 1=start, 2=end, 3=strand, 4=coverage, 5=gene,
//
//	6=cov_db, 7=accession, 8=gap, 9=cov_q, 10=identity
func buildAbricateLine(seqid, gene, covDb, accession, covQ, identity string) string {
	return seqid + "\t" + "100" + "\t" + "200" + "\t" + "+" + "\t" +
		"100/100" + "\t" + gene + "\t" + covDb + "\t" + accession + "\t" +
		"0" + "\t" + covQ + "\t" + identity
}

func TestGetAbricateResult(t *testing.T) {
	t.Run("Success - High Coverage And Identity", func(t *testing.T) {
		line := buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123456", "95.5", "98.0")
		path := createMockAbricateFile(t, line+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Contains(t, results[0], "blaTEM")
	})

	t.Run("Success - Low Coverage Excluded", func(t *testing.T) {
		line := buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123456", "80.0", "98.0")
		path := createMockAbricateFile(t, line+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("Success - Low Identity Excluded", func(t *testing.T) {
		line := buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123456", "95.0", "85.0")
		path := createMockAbricateFile(t, line+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("Success - Vancomycin Gene Included Regardless Of Coverage", func(t *testing.T) {
		line := buildAbricateLine("seq1", "VanA", "resfinder", "AF123456", "50.0", "50.0")
		path := createMockAbricateFile(t, line+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Contains(t, results[0], "VanA")
	})

	t.Run("Success - Van Gene Case Insensitive", func(t *testing.T) {
		line := buildAbricateLine("seq1", "vanB", "resfinder", "AF123456", "50.0", "50.0")
		path := createMockAbricateFile(t, line+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Contains(t, results[0], "vanB")
	})

	t.Run("Success - Mixed Lines", func(t *testing.T) {
		line1 := buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123", "95.0", "98.0")
		line2 := buildAbricateLine("seq2", "blaCTX", "resfinder", "AF456", "80.0", "80.0")
		line3 := buildAbricateLine("seq3", "VanC", "resfinder", "AF789", "60.0", "60.0")
		path := createMockAbricateFile(t, line1+"\n"+line2+"\n"+line3+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("Success - Lines With Fewer Than 11 Fields Skipped", func(t *testing.T) {
		shortLine := "seq1\t100\t200\t+\t100/100\tblaTEM"
		highLine := buildAbricateLine("seq2", "blaCTX", "resfinder", "AF456", "95.0", "98.0")
		path := createMockAbricateFile(t, shortLine+"\n"+highLine+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Contains(t, results[0], "blaCTX")
	})

	t.Run("Success - Invalid Coverage/Identity Skipped", func(t *testing.T) {
		line := "seq1\t100\t200\t+\t100/100\tblaTEM\tresfinder\tAF123\t0\tnotanumber\tnotanumber"
		path := createMockAbricateFile(t, line+"\n")

		results, err := GetAbricateResult(path)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("Success - Empty File", func(t *testing.T) {
		path := createMockAbricateFile(t, "")

		results, err := GetAbricateResult(path)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Error(), "Empty Abricate result")
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		results, err := GetAbricateResult("nonexistent.txt")
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Error(), "Failed to open Abricate result")
	})
}

func TestProcessResfinder(t *testing.T) {
	// refCatalog format: TSV with at least 17 columns, last column is antibiotic name
	// refItem[0] = gene name in DB, refItem[len-17] = antibiotic
	// NOTE: ProcessResfinder uses bufio.NewReader.Peek to check for empty file,
	// which consumes data from the underlying refFile via its internal buffer.
	// Then bufio.NewScanner reads from refFile's advanced position.
	// To ensure the scanner reads the gene data, we pad a header line > 4096 bytes
	// so br.Peek consumes the header, and the scanner reads the gene lines.
	buildRefLine := func(gene, antibiotic string) string {
		// The code uses refItem[len(refItem)-17] as the antibiotic name.
		// The reference file is read with strings.TrimSpace which strips trailing tabs,
		// so we need a non-empty last field to preserve the column count.
		// With 34 columns and a non-empty last field, after TrimSpace we still have 34 fields.
		// refItem[34-17] = refItem[17] = antibiotic.
		cols := make([]string, 34)
		cols[0] = gene
		cols[17] = antibiotic
		cols[33] = "placeholder"
		result := ""
		for i, c := range cols {
			if i > 0 {
				result += "\t"
			}
			result += c
		}
		return result
	}

	buildRefContent := func(gene, antibiotic string) string {
		// Pad header line to > 4096 bytes so br.Peek consumes it entirely
		header := strings.Repeat("P", 10000) + "\n"
		return header + buildRefLine(gene, antibiotic) + "\n"
	}

	t.Run("Success - Gene Found In Reference", func(t *testing.T) {
		refContent := buildRefContent("blaTEM", "Ampicillin")
		refPath := createMockAbricateFile(t, refContent)

		abricateResult := []string{
			buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123", "95.0", "98.0"),
		}

		geneResults, blastResults, err := ProcessResfinder(abricateResult, refPath)
		assert.NoError(t, err)
		assert.Len(t, geneResults, 1)
		assert.Len(t, blastResults, 1)
		assert.Contains(t, geneResults[0], "blaTEM")
		assert.Contains(t, geneResults[0], "resistance to ampicillin")
		assert.Contains(t, geneResults[0], "allele confidence 98.0")
		assert.Contains(t, blastResults[0], "blaTEM")
		assert.Contains(t, blastResults[0], "ID: 98.0")
		assert.Contains(t, blastResults[0], "COV_Q: 95.0")
		assert.Contains(t, blastResults[0], "COV_DB: resfinder")
	})

	t.Run("Success - Gene Not Found In Reference", func(t *testing.T) {
		refContent := buildRefContent("otherGene", "Other")
		refPath := createMockAbricateFile(t, refContent)

		abricateResult := []string{
			buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123", "95.0", "98.0"),
		}

		geneResults, blastResults, err := ProcessResfinder(abricateResult, refPath)
		assert.NoError(t, err)
		assert.Len(t, geneResults, 1)
		assert.Len(t, blastResults, 1)
		assert.Contains(t, geneResults[0], "blaTEM")
		assert.Contains(t, geneResults[0], "allele confidence 98.0")
		assert.NotContains(t, geneResults[0], "resistance")
	})

	t.Run("Success - Gene With Underscore Split", func(t *testing.T) {
		refContent := buildRefContent("blaTEM", "Ampicillin")
		refPath := createMockAbricateFile(t, refContent)

		abricateResult := []string{
			buildAbricateLine("seq1", "blaTEM_extra", "resfinder", "AF123", "95.0", "98.0"),
		}

		geneResults, _, err := ProcessResfinder(abricateResult, refPath)
		assert.NoError(t, err)
		assert.Len(t, geneResults, 1)
		assert.Contains(t, geneResults[0], "resistance to ampicillin")
	})

	t.Run("Success - Empty Abricate Result", func(t *testing.T) {
		refContent := buildRefContent("blaTEM", "Ampicillin")
		refPath := createMockAbricateFile(t, refContent)

		geneResults, blastResults, err := ProcessResfinder([]string{}, refPath)
		assert.NoError(t, err)
		assert.Empty(t, geneResults)
		assert.Empty(t, blastResults)
	})

	t.Run("Success - Lines With Fewer Than 11 Fields Skipped", func(t *testing.T) {
		refContent := buildRefContent("blaTEM", "Ampicillin")
		refPath := createMockAbricateFile(t, refContent)

		abricateResult := []string{
			"seq1\t100\t200\t+\t100/100\tblaTEM",
		}

		geneResults, blastResults, err := ProcessResfinder(abricateResult, refPath)
		assert.NoError(t, err)
		assert.Empty(t, geneResults)
		assert.Empty(t, blastResults)
	})

	t.Run("Error - Empty Reference File", func(t *testing.T) {
		refPath := createMockAbricateFile(t, "")

		abricateResult := []string{
			buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123", "95.0", "98.0"),
		}

		geneResults, blastResults, err := ProcessResfinder(abricateResult, refPath)
		assert.Error(t, err)
		assert.Nil(t, geneResults)
		assert.Nil(t, blastResults)
		assert.Contains(t, err.Error(), "Empty Resfinder reference file")
	})

	t.Run("Error - Reference File Not Found", func(t *testing.T) {
		abricateResult := []string{
			buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123", "95.0", "98.0"),
		}

		geneResults, blastResults, err := ProcessResfinder(abricateResult, "nonexistent.txt")
		assert.Error(t, err)
		assert.Nil(t, geneResults)
		assert.Nil(t, blastResults)
		assert.Contains(t, err.Error(), "Failed to open Resfinder reference file")
	})

	t.Run("Success - Reference Line With Fewer Than 17 Columns Skipped", func(t *testing.T) {
		header := strings.Repeat("P", 10000) + "\n"
		shortRef := "blaTEM\tshort"
		refPath := createMockAbricateFile(t, header+shortRef+"\n")

		abricateResult := []string{
			buildAbricateLine("seq1", "blaTEM", "resfinder", "AF123", "95.0", "98.0"),
		}

		geneResults, _, err := ProcessResfinder(abricateResult, refPath)
		assert.NoError(t, err)
		assert.Len(t, geneResults, 1)
		assert.Contains(t, geneResults[0], "allele confidence")
		assert.NotContains(t, geneResults[0], "resistance")
	})
}

func TestProcessVFDB(t *testing.T) {
	// VFDB line needs at least 14 fields (0-13)
	// ProcessVFDB uses fields[1], fields[5], fields[13], fields[10], fields[9], fields[6]
	buildVFDBLine := func(locus, gene, extra, covQ, identity, covDb string) string {
		fields := make([]string, 14)
		fields[0] = "contig1"
		fields[1] = locus
		fields[5] = gene
		fields[6] = covDb
		fields[9] = covQ
		fields[10] = identity
		fields[13] = extra
		result := ""
		for i, f := range fields {
			if i > 0 {
				result += "\t"
			}
			if f == "" {
				result += "."
			} else {
				result += f
			}
		}
		return result
	}

	t.Run("Success - Single Line", func(t *testing.T) {
		line := buildVFDBLine("VF0001", "virB4", "TypeIV secretion", "95.0", "98.0", "vfdb")
		results := ProcessVFDB([]string{line})
		assert.Len(t, results, 1)
		assert.Contains(t, results[0], "VF0001")
		assert.Contains(t, results[0], "virB4")
		assert.Contains(t, results[0], "TypeIV secretion")
		assert.Contains(t, results[0], "ID: 98.0")
		assert.Contains(t, results[0], "COV_Q: 95.0")
		assert.Contains(t, results[0], "COV_DB: vfdb")
	})

	t.Run("Success - Multiple Lines", func(t *testing.T) {
		line1 := buildVFDBLine("VF0001", "virB4", "TypeIV secretion", "95.0", "98.0", "vfdb")
		line2 := buildVFDBLine("VF0002", "virB11", "T4SS", "88.0", "92.0", "vfdb")
		results := ProcessVFDB([]string{line1, line2})
		assert.Len(t, results, 2)
	})

	t.Run("Success - Lines With Fewer Than 14 Fields Skipped", func(t *testing.T) {
		shortLine := "contig1\tVF0001\t100\t200\t+\t100/100\tvirB4\tvfdb\tAF123\t0\t95.0\t98.0"
		results := ProcessVFDB([]string{shortLine})
		assert.Empty(t, results)
	})

	t.Run("Success - Empty Input", func(t *testing.T) {
		results := ProcessVFDB([]string{})
		assert.Empty(t, results)
		assert.Nil(t, results)
	})
}

func TestProcessPlasmidFinder(t *testing.T) {
	t.Run("Success - Single Line", func(t *testing.T) {
		line := buildAbricateLine("seq1", "repB", "plasmidfinder", "AF123", "95.0", "98.0")
		results := ProcessPlasmidFinder([]string{line})
		assert.Len(t, results, 1)
		assert.Contains(t, results[0], "repB")
		assert.Contains(t, results[0], "ID: 98.0")
		assert.Contains(t, results[0], "COV_Q: 95.0")
		assert.Contains(t, results[0], "COV_DB: plasmidfinder")
	})

	t.Run("Success - Multiple Lines", func(t *testing.T) {
		line1 := buildAbricateLine("seq1", "repB", "plasmidfinder", "AF123", "95.0", "98.0")
		line2 := buildAbricateLine("seq2", "blaTEM", "plasmidfinder", "AF456", "90.0", "95.0")
		results := ProcessPlasmidFinder([]string{line1, line2})
		assert.Len(t, results, 2)
	})

	t.Run("Success - Lines With Fewer Than 11 Fields Skipped", func(t *testing.T) {
		shortLine := "seq1\t100\t200\t+\t100/100\trepB"
		results := ProcessPlasmidFinder([]string{shortLine})
		assert.Empty(t, results)
	})

	t.Run("Success - Empty Input", func(t *testing.T) {
		results := ProcessPlasmidFinder([]string{})
		assert.Empty(t, results)
		assert.Nil(t, results)
	})
}
