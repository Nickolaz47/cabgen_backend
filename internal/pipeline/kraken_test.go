package pipeline

import (
	"fmt"
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

func krakenReportLine(name string, cladeReads int) string {
	return fmt.Sprintf("1.00\t%d\t%d\tS\t0\t%s\n", cladeReads, cladeReads,
		name)
}

func TestKrakenSpeciesCounter(t *testing.T) {
	t.Run("Success - Single Clade Each Species", func(t *testing.T) {
		mockContent := krakenReportLine("Escherichia coli", 1) +
			krakenReportLine("Klebsiella pneumoniae", 1)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.NotNil(t, second)

		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Equal(t, "Klebsiella pneumoniae", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Multiple Clades Same Species", func(t *testing.T) {
		mockContent := krakenReportLine("Escherichia coli", 3) +
			krakenReportLine("Klebsiella pneumoniae", 1)
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

	t.Run("Success - Species Name Kept As-Is", func(t *testing.T) {
		mockContent := krakenReportLine("Escherichia coli (strain K12)", 2) +
			krakenReportLine("Klebsiella pneumoniae", 1)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli (strain K12)", first.Name)
		assert.Equal(t, 2, first.Count)
		assert.NotNil(t, second)
		assert.Equal(t, "Klebsiella pneumoniae", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Species Name With Indentation", func(t *testing.T) {
		mockContent := krakenReportLine("  Escherichia coli", 2) +
			krakenReportLine("Klebsiella pneumoniae", 1)
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
		mockContent := krakenReportLine("Zebra", 1) +
			krakenReportLine("Alpha", 1)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Alpha", first.Name)
		assert.NotNil(t, second)
		assert.Equal(t, "Zebra", second.Name)
	})

	t.Run("Success - Single Species Returns Nil Second", func(t *testing.T) {
		mockContent := krakenReportLine("Escherichia coli", 2)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 2, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Non-Species Ranks Ignored", func(t *testing.T) {
		mockContent := "50.00\t10\t10\tG\t561\tEscherichia" + "\n" +
			krakenReportLine("Escherichia coli", 1)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Lines With Fewer Than 6 Fields Ignored", func(t *testing.T) {
		mockContent := "1.00\t1\t1\tS\n" +
			krakenReportLine("Escherichia coli", 1)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Escherichia coli", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Empty Species Name Ignored", func(t *testing.T) {
		mockContent := krakenReportLine("Escherichia coli", 1) +
			"1.00\t5\t5\tS\t0\t\n" +
			krakenReportLine("Klebsiella pneumoniae", 1)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.NotNil(t, second)

		names := []string{first.Name, second.Name}
		assert.ElementsMatch(t,
			[]string{"Escherichia coli", "Klebsiella pneumoniae"}, names)
		assert.Equal(t, 1, first.Count)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Invalid Clade Reads Ignored", func(t *testing.T) {
		mockContent := "1.00\tnot-a-number\t1\tS\t0\tEscherichia coli\n" +
			krakenReportLine("Klebsiella pneumoniae", 1)
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "Klebsiella pneumoniae", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Empty File Returns Error", func(t *testing.T) {
		path := createMockKrakenFile(t, "")

		first, second, err := KrakenSpeciesCounter(path)
		assert.Error(t, err)
		assert.Nil(t, first)
		assert.Nil(t, second)
		assert.ErrorContains(t, err, "Empty Kraken report")
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		first, second, err := KrakenSpeciesCounter("nonexistent_path.txt")
		assert.Error(t, err)
		assert.Nil(t, first)
		assert.Nil(t, second)
		assert.Contains(t, err.Error(), "Kraken report file not found")
	})
}
