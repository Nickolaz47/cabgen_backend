package pipeline

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockKrakenFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "kraken_mock_*.txt")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

// krakenLine builds a kraken2 classified line: "C\tseq_id\ttaxon\tkmers".
func krakenLine(seqID, taxon string) string {
	return "C\t" + seqID + "\t" + taxon + "\t|0:0|\n"
}

func TestKrakenSpeciesCounter(t *testing.T) {
	t.Run("Success - Single Read Each Species", func(t *testing.T) {
		mockContent := krakenLine("read1", "Escherichia coli") +
			krakenLine("read2", "Klebsiella pneumoniae")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.NotNil(t, second)

		// tie broken by name ascending: "Escherichia coli" < "Klebsiella pneumoniae"
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Equal(t, "Klebsiella pneumoniae", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Multiple Reads Same Species", func(t *testing.T) {
		mockContent := krakenLine("read1", "Escherichia coli") +
			krakenLine("read2", "Escherichia coli") +
			krakenLine("read3", "Escherichia coli") +
			krakenLine("read4", "Klebsiella pneumoniae")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 3, first.Count)
		assert.NotNil(t, second)
		assert.Equal(t, "Klebsiella pneumoniae", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Species With Parentheses", func(t *testing.T) {
		mockContent := krakenLine("read1", "Escherichia coli (strain K12)") +
			krakenLine("read2", "Escherichia coli (strain B)") +
			krakenLine("read3", "Klebsiella pneumoniae")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 2, first.Count)
		assert.NotNil(t, second)
		assert.Equal(t, "Klebsiella pneumoniae", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Tie Broken By Name Ascending", func(t *testing.T) {
		mockContent := krakenLine("read1", "Zebra") +
			krakenLine("read2", "Alpha")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Alpha", first.Name)
		assert.NotNil(t, second)
		assert.Equal(t, "Zebra", second.Name)
	})

	t.Run("Success - Single Species Returns Nil Second", func(t *testing.T) {
		mockContent := krakenLine("read1", "Escherichia coli") +
			krakenLine("read2", "Escherichia coli")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 2, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Unclassified Lines Ignored", func(t *testing.T) {
		mockContent := "U\tread1\t0\t|0:0|" + "\n" +
			krakenLine("read2", "Escherichia coli")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Lines With Fewer Than 3 Fields Ignored", func(t *testing.T) {
		mockContent := "C\tread1\n" +
			krakenLine("read2", "Escherichia coli")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Empty Species Name Ignored", func(t *testing.T) {
		mockContent := krakenLine("read1", "Escherichia coli") +
			"C\tread2\t\t|0:0|\n" +
			krakenLine("read3", "Klebsiella pneumoniae")
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.NotNil(t, second)

		names := []string{first.Name, second.Name}
		assert.ElementsMatch(t, []string{"Escherichia coli", "Klebsiella pneumoniae"}, names)
		assert.Equal(t, 1, first.Count)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Empty File Returns Error", func(t *testing.T) {
		path := createMockKrakenFile(t, "")

		first, second, err := KrakenSpeciesCounter(path)
		assert.Error(t, err)
		assert.Nil(t, first)
		assert.Nil(t, second)
		assert.ErrorContains(t, err, "Empty Kraken result")
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		first, second, err := KrakenSpeciesCounter("nonexistent_path.txt")
		assert.Error(t, err)
		assert.Nil(t, first)
		assert.Nil(t, second)
		assert.Contains(t, err.Error(), "Kraken output file not found")
	})
}