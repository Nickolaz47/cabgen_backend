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

func TestKrakenSpeciesCounter(t *testing.T) {
	t.Run("Success - Single Species Each", func(t *testing.T) {
		mockContent := "read1\t0\tSpeciesA\nread2\t1\tSpeciesB\n"
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "SpeciesA", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.NotNil(t, second)
		assert.Equal(t, "SpeciesB", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Multiple Reads Same Species", func(t *testing.T) {
		mockContent := "read1\t0\tSpeciesA\nread2\t1\tSpeciesA\nread3\t2\tSpeciesA\nread4\t3\tSpeciesB\n"
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "SpeciesA", first.Name)
		assert.Equal(t, 3, first.Count)
		assert.NotNil(t, second)
		assert.Equal(t, "SpeciesB", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Species With Parentheses", func(t *testing.T) {
		mockContent := "read1\t0\tSpeciesA (strain X)\nread2\t1\tSpeciesA (strain Y)\nread3\t2\tSpeciesB\n"
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "SpeciesA", first.Name)
		assert.Equal(t, 2, first.Count)
		assert.NotNil(t, second)
		assert.Equal(t, "SpeciesB", second.Name)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Success - Single Species Returns Nil Second", func(t *testing.T) {
		mockContent := "read1\t0\tSpeciesA\nread2\t1\tSpeciesA\n"
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "SpeciesA", first.Name)
		assert.Equal(t, 2, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Lines With Fewer Than 3 Fields Ignored", func(t *testing.T) {
		mockContent := "read1\t0\nread2\t1\tSpeciesA\n"
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.Equal(t, "SpeciesA", first.Name)
		assert.Equal(t, 1, first.Count)
		assert.Nil(t, second)
	})

	t.Run("Success - Empty Species Name Ignored", func(t *testing.T) {
		mockContent := "read1\t0\tSpeciesA\nread2\t1\t\nread3\t2\tSpeciesB\n"
		path := createMockKrakenFile(t, mockContent)

		first, second, err := KrakenSpeciesCounter(path)
		assert.NoError(t, err)
		assert.NotNil(t, first)
		assert.NotNil(t, second)

		names := []string{first.Name, second.Name}
		assert.ElementsMatch(t, []string{"SpeciesA", "SpeciesB"}, names)
		assert.Equal(t, 1, first.Count)
		assert.Equal(t, 1, second.Count)
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		first, second, err := KrakenSpeciesCounter("nonexistent_path.txt")
		assert.Error(t, err)
		assert.Nil(t, first)
		assert.Nil(t, second)
		assert.Contains(t, err.Error(), "Kraken output file not found")
	})

	t.Run("Error - Empty File", func(t *testing.T) {
		path := createMockKrakenFile(t, "")

		first, second, err := KrakenSpeciesCounter(path)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Empty Kraken result")
		assert.Nil(t, first)
		assert.Nil(t, second)
	})
}
