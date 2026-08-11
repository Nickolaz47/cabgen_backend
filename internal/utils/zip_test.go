package utils_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestZipDirectory(t *testing.T) {
	t.Run("Success - Zips Files and Subdirectories", func(t *testing.T) {
		srcDir := t.TempDir()
		assert.NoError(t, os.MkdirAll(filepath.Join(srcDir, "report"), 0755))
		assert.NoError(t, os.MkdirAll(filepath.Join(srcDir, "qc"), 0755))
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "report", "summary.json"),
			[]byte(`{"ok":true}`), 0644))
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "qc", "reads_fastqc.html"),
			[]byte("<html></html>"), 0644))
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "already.zip"),
			[]byte("x"), 0644))

		destPath := filepath.Join(t.TempDir(), "results.zip")
		err := utils.ZipDirectory(srcDir, destPath)

		assert.NoError(t, err)

		zr, err := zip.OpenReader(destPath)
		assert.NoError(t, err)
		defer zr.Close()

		names := make(map[string]bool)
		for _, f := range zr.File {
			names[f.Name] = true
		}

		assert.True(t, names[filepath.Join(filepath.Base(srcDir),
			"report", "summary.json")])
		assert.True(t, names[filepath.Join(filepath.Base(srcDir),
			"qc", "reads_fastqc.html")])
		assert.False(t, names[filepath.Join(filepath.Base(srcDir),
			"already.zip")],
			"archive files must be excluded from the zip")
	})

	t.Run("Error - Invalid Source Directory", func(t *testing.T) {
		err := utils.ZipDirectory(
			"/nonexistent/src", "/nonexistent/out.zip")

		assert.Error(t, err)
	})
}

func TestZipSubdirectories(t *testing.T) {
	t.Run("Success - Zips Only Selected Subdirectories", func(t *testing.T) {
		srcDir := t.TempDir()
		for _, dir := range []string{"qc", "assembly", "amr", "report"} {
			assert.NoError(t, os.MkdirAll(filepath.Join(srcDir, dir),
				0755))
		}
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "qc", "reads_fastqc.html"),
			[]byte("<html></html>"), 0644))
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "assembly", "genome.fasta"),
			[]byte(">genome"), 0644))
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "amr", "resfinder.tsv"),
			[]byte("gene"), 0644))
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "qc", "existing_results.zip"),
			[]byte("x"), 0644))
		assert.NoError(t, os.WriteFile(
			filepath.Join(srcDir, "report", "results.zip"),
			[]byte("x"), 0644))

		destPath := filepath.Join(t.TempDir(), "out.zip")
		err := utils.ZipSubdirectories(srcDir,
			[]string{"qc", "assembly", "amr"}, destPath)

		assert.NoError(t, err)

		zr, err := zip.OpenReader(destPath)
		assert.NoError(t, err)
		defer zr.Close()

		var names []string
		for _, f := range zr.File {
			names = append(names, f.Name)
		}
		base := filepath.Base(srcDir)
		assert.Contains(t, names, filepath.Join(base, "qc",
			"reads_fastqc.html"))
		assert.Contains(t, names, filepath.Join(base, "assembly",
			"genome.fasta"))
		assert.Contains(t, names, filepath.Join(base, "amr",
			"resfinder.tsv"))
		assert.NotContains(t, names, filepath.Join(base, "report",
			"results.zip"))
		assert.NotContains(t, names, filepath.Join(base, "qc",
			"existing_results.zip"))
	})

	t.Run("Error - Missing Subdirectory", func(t *testing.T) {
		srcDir := t.TempDir()

		err := utils.ZipSubdirectories(srcDir,
			[]string{"qc"}, filepath.Join(t.TempDir(), "out.zip"))

		assert.Error(t, err)
	})
}
