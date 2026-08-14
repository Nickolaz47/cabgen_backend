package pipeline

import (
	"bytes"
	"compress/gzip"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockFastqFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "mock_*.fastq")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

func createMockGzipFastqFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "mock_*.fastq.gz")
	assert.NoError(t, err)

	gzWriter := gzip.NewWriter(tmpFile)
	_, err = gzWriter.Write([]byte(content))
	assert.NoError(t, err)

	err = gzWriter.Close()
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

func fastqRead(id, seq string) string {
	return id + "\n" + seq + "\n+\n" + qualityScores(len(seq)) + "\n"
}

func qualityScores(n int) string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.WriteByte('I')
	}
	return buf.String()
}

func TestProcessFastq(t *testing.T) {
	t.Run("Success - Single Read", func(t *testing.T) {
		content := fastqRead("@read1", "ATCGATCG")
		path := createMockFastqFile(t, content)

		readCount, totalBases, err := processFastq(path)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), readCount)
		assert.Equal(t, int64(8), totalBases)
	})

	t.Run("Success - Multiple Reads", func(t *testing.T) {
		content := fastqRead("@read1", "ATCGATCG") + fastqRead("@read2", "GCTAGCTA")
		path := createMockFastqFile(t, content)

		readCount, totalBases, err := processFastq(path)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), readCount)
		assert.Equal(t, int64(16), totalBases)
	})

	t.Run("Success - Reads With Different Lengths", func(t *testing.T) {
		content := fastqRead("@read1", "ATCG") + fastqRead("@read2", "GCTAGCTAAA")
		path := createMockFastqFile(t, content)

		readCount, totalBases, err := processFastq(path)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), readCount)
		assert.Equal(t, int64(14), totalBases)
	})

	t.Run("Success - Gzipped Fastq", func(t *testing.T) {
		content := fastqRead("@read1", "ATCGATCG") + fastqRead("@read2", "GCTAGCTA")
		path := createMockGzipFastqFile(t, content)

		readCount, totalBases, err := processFastq(path)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), readCount)
		assert.Equal(t, int64(16), totalBases)
	})

	t.Run("Error - Empty File Returns Empty Reads", func(t *testing.T) {
		path := createMockFastqFile(t, "")

		readCount, totalBases, err := processFastq(path)
		assert.ErrorIs(t, err, ErrEmptyReads)
		assert.Equal(t, int64(0), readCount)
		assert.Equal(t, int64(0), totalBases)
	})

	t.Run("Error - FASTA File Returns Invalid Format", func(t *testing.T) {
		path := createMockFastqFile(t, ">seq1\nATCGATCG\n")

		readCount, totalBases, err := processFastq(path)
		assert.ErrorIs(t, err, ErrInvalidFormat)
		assert.Equal(t, int64(0), readCount)
		assert.Equal(t, int64(0), totalBases)
	})

	t.Run("Error - Invalid Header Returns Invalid Format", func(t *testing.T) {
		path := createMockFastqFile(t, "not_a_header\nATCG\n")

		readCount, totalBases, err := processFastq(path)
		assert.ErrorIs(t, err, ErrInvalidFormat)
		assert.Equal(t, int64(0), readCount)
		assert.Equal(t, int64(0), totalBases)
	})

	t.Run("Error - Truncated FASTQ Returns Corrupted Input", func(t *testing.T) {
		path := createMockFastqFile(t, "@read1\nATCGATCG\n+\n")

		readCount, totalBases, err := processFastq(path)
		assert.ErrorIs(t, err, ErrCorruptedInput)
		assert.Equal(t, int64(0), readCount)
		assert.Equal(t, int64(0), totalBases)
	})

	t.Run("Success - Sequence With Whitespace Trimmed", func(t *testing.T) {
		content := "@read1\n  ATCGATCG  \n+\nIIIIIIII\n"
		path := createMockFastqFile(t, content)

		readCount, totalBases, err := processFastq(path)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), readCount)
		assert.Equal(t, int64(8), totalBases)
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		readCount, totalBases, err := processFastq("nonexistent.fastq")
		assert.ErrorIs(t, err, ErrFileNotFound)
		assert.Equal(t, int64(0), readCount)
		assert.Equal(t, int64(0), totalBases)
	})
}

func TestCalculateCoverage(t *testing.T) {
	t.Run("Success - Calculates Coverage", func(t *testing.T) {
		// 4 reads, avg len 8, genome 100 → (8*4)/100 = 0.32
		read1 := createMockFastqFile(t, fastqRead("@r1", "ATCGATCG")+fastqRead("@r2", "ATCGATCG"))
		read2 := createMockFastqFile(t, fastqRead("@r3", "GCTAGCTA")+fastqRead("@r4", "GCTAGCTA"))

		coverage, err := CalculateCoverage(read1, read2, 100)
		assert.NoError(t, err)
		assert.Equal(t, 0.32, coverage)
	})

	t.Run("Success - Single Read Per File", func(t *testing.T) {
		// 2 reads, avg len 10, genome 50 → (10*2)/50 = 0.4
		read1 := createMockFastqFile(t, fastqRead("@r1", "ATCGATCGAT"))
		read2 := createMockFastqFile(t, fastqRead("@r2", "GCTAGCTAGC"))

		coverage, err := CalculateCoverage(read1, read2, 50)
		assert.NoError(t, err)
		assert.Equal(t, 0.4, coverage)
	})

	t.Run("Success - Gzipped Files", func(t *testing.T) {
		read1 := createMockGzipFastqFile(t, fastqRead("@r1", "ATCGATCG")+fastqRead("@r2", "ATCGATCG"))
		read2 := createMockGzipFastqFile(t, fastqRead("@r3", "GCTAGCTA")+fastqRead("@r4", "GCTAGCTA"))

		coverage, err := CalculateCoverage(read1, read2, 100)
		assert.NoError(t, err)
		assert.Equal(t, 0.32, coverage)
	})

	t.Run("Success - Reads With Different Lengths", func(t *testing.T) {
		// 3 reads, avg len 5, genome 100 → (5*3)/100 = 0.15
		read1 := createMockFastqFile(t, fastqRead("@r1", "ATCG")+fastqRead("@r2", "ATCGAT"))
		read2 := createMockFastqFile(t, fastqRead("@r3", "GCTAGCTA"))

		coverage, err := CalculateCoverage(read1, read2, 100)
		assert.NoError(t, err)
		assert.Equal(t, 0.15, coverage)
	})

	t.Run("Error - Genome Size Zero", func(t *testing.T) {
		read1 := createMockFastqFile(t, fastqRead("@r1", "ATCG"))
		read2 := createMockFastqFile(t, fastqRead("@r2", "GCTA"))

		coverage, err := CalculateCoverage(read1, read2, 0)
		assert.Error(t, err)
		assert.Equal(t, float64(0), coverage)
		assert.Contains(t, err.Error(), "Invalid genome size")
	})

	t.Run("Error - Negative Genome Size", func(t *testing.T) {
		read1 := createMockFastqFile(t, fastqRead("@r1", "ATCG"))
		read2 := createMockFastqFile(t, fastqRead("@r2", "GCTA"))

		coverage, err := CalculateCoverage(read1, read2, -100)
		assert.Error(t, err)
		assert.Equal(t, float64(0), coverage)
		assert.Contains(t, err.Error(), "Invalid genome size")
	})

	t.Run("Error - Read1 File Not Found", func(t *testing.T) {
		read2 := createMockFastqFile(t, fastqRead("@r2", "GCTA"))

		coverage, err := CalculateCoverage("nonexistent.fastq", read2, 100)
		assert.ErrorIs(t, err, ErrFileNotFound)
		assert.Equal(t, float64(0), coverage)
		assert.Contains(t, err.Error(), "Failed to process read1")
	})

	t.Run("Error - Read2 File Not Found", func(t *testing.T) {
		read1 := createMockFastqFile(t, fastqRead("@r1", "ATCG"))

		coverage, err := CalculateCoverage(read1, "nonexistent.fastq", 100)
		assert.ErrorIs(t, err, ErrFileNotFound)
		assert.Equal(t, float64(0), coverage)
		assert.Contains(t, err.Error(), "Failed to process read2")
	})

	t.Run("Error - Empty Read1 Returns Empty Reads", func(t *testing.T) {
		read1 := createMockFastqFile(t, "")
		read2 := createMockFastqFile(t, fastqRead("@r2", "GCTA"))

		coverage, err := CalculateCoverage(read1, read2, 100)
		assert.ErrorIs(t, err, ErrEmptyReads)
		assert.Equal(t, float64(0), coverage)
		assert.Contains(t, err.Error(), "Failed to process read1")
	})

	t.Run("Error - FASTA Read1 Returns Invalid Format", func(t *testing.T) {
		read1 := createMockFastqFile(t, ">seq1\nATCG\n")
		read2 := createMockFastqFile(t, fastqRead("@r2", "GCTA"))

		coverage, err := CalculateCoverage(read1, read2, 100)
		assert.ErrorIs(t, err, ErrInvalidFormat)
		assert.Equal(t, float64(0), coverage)
	})

	t.Run("Error - Truncated Read1 Returns Corrupted Input", func(t *testing.T) {
		read1 := createMockFastqFile(t, "@r1\nATCG\n+\n")
		read2 := createMockFastqFile(t, fastqRead("@r2", "GCTA"))

		coverage, err := CalculateCoverage(read1, read2, 100)
		assert.ErrorIs(t, err, ErrCorruptedInput)
		assert.Equal(t, float64(0), coverage)
	})
}
