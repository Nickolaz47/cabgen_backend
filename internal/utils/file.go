package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// ResolveSampleFilePath resolves the absolute path for a file in a sample's
// upload directory. For FastQ files, it returns the path without checking
// existence. For FASTA files, it first checks the sample directory, then
// falls back to the analysis assembly directory.
func ResolveSampleFilePath(rootDir, userID, sampleID, fileName, fileType string, analysisID string) (string, bool) {
	sampleDir := filepath.Join(rootDir, "uploads", "users", userID, "samples", sampleID)

	if fileType != "fasta" {
		return filepath.Join(sampleDir, fileName), true
	}

	// FASTA: try sample directory first
	samplePath := filepath.Join(sampleDir, fileName)
	if _, err := os.Stat(samplePath); err == nil {
		return samplePath, true
	}

	// Fallback: analysis assembly directory
	if analysisID != "" {
		assemblyPath := filepath.Join(sampleDir, "analyses", analysisID, "assembly", fileName)
		if _, err := os.Stat(assemblyPath); err == nil {
			return assemblyPath, true
		}
	}

	return "", false
}
