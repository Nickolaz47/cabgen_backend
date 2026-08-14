package pipeline

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func processFastq(filePath string) (int64, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, 0, ErrFileNotFound
		}
		return 0, 0, fmt.Errorf("%w: %v", ErrCorruptedInput, err)
	}
	defer file.Close()

	var reader io.Reader = file

	buf := make([]byte, 2)
	if _, err := file.Read(buf); err == nil {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return 0, 0, fmt.Errorf("%w: %v", ErrCorruptedInput, err)
		}

		if buf[0] == 0x1f && buf[1] == 0x8b {
			gzReader, err := gzip.NewReader(file)
			if err != nil {
				return 0, 0, fmt.Errorf("%w: %v", ErrCorruptedInput, err)
			}
			defer gzReader.Close()
			reader = gzReader
		}
	} else if errors.Is(err, io.EOF) {
		return 0, 0, ErrEmptyReads
	} else {
		return 0, 0, fmt.Errorf("%w: %v", ErrCorruptedInput, err)
	}

	scanner := bufio.NewScanner(reader)

	const maxCapacity = 1024 * 1024 * 4
	lineBuf := make([]byte, maxCapacity)
	scanner.Buffer(lineBuf, maxCapacity)

	var lineCount int64 = 0
	var readCount int64 = 0
	var totalBases int64 = 0

	for scanner.Scan() {
		lineCount++

		if lineCount == 1 {
			line := strings.TrimSpace(scanner.Text())
			switch {
			case strings.HasPrefix(line, ">"):
				return 0, 0, ErrInvalidFormat
			case !strings.HasPrefix(line, "@"):
				return 0, 0, ErrInvalidFormat
			}
		}

		if lineCount%4 == 2 {
			readCount++
			totalBases += int64(len(strings.TrimSpace(scanner.Text())))
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("%w: %v", ErrCorruptedInput, err)
	}

	if lineCount%4 != 0 {
		return 0, 0, ErrCorruptedInput
	}

	return readCount, totalBases, nil
}

func CalculateCoverage(read1, read2 string, genomeSize int64) (float64, error) {
	if genomeSize <= 0 {
		return 0, fmt.Errorf("Invalid genome size: %d", genomeSize)
	}

	count1, bases1, err := processFastq(read1)
	if err != nil {
		return 0, fmt.Errorf("Failed to process read1: %w", err)
	}

	count2, _, err := processFastq(read2)
	if err != nil {
		return 0, fmt.Errorf("Failed to process read2: %w", err)
	}

	if count1 == 0 {
		return 0, ErrEmptyReads
	}

	totalReads := float64(count1 + count2)
	avgLength := float64(bases1) / float64(count1)
	coverage := (avgLength * totalReads) / float64(genomeSize)

	return coverage, nil
}
