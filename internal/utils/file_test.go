package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tmpDir := t.TempDir()
		src := filepath.Join(tmpDir, "source.fasta")
		dst := filepath.Join(tmpDir, "dest.fasta")

		content := ">seq1\nATCGATCG\n"
		err := os.WriteFile(src, []byte(content), 0644)
		assert.NoError(t, err)

		err = CopyFile(src, dst)
		assert.NoError(t, err)

		result, err := os.ReadFile(dst)
		assert.NoError(t, err)
		assert.Equal(t, content, string(result))
	})

	t.Run("Error - Source Not Found", func(t *testing.T) {
		tmpDir := t.TempDir()
		src := filepath.Join(tmpDir, "nonexistent.fasta")
		dst := filepath.Join(tmpDir, "dest.fasta")

		err := CopyFile(src, dst)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open source file")
	})

	t.Run("Error - Invalid Destination", func(t *testing.T) {
		tmpDir := t.TempDir()
		src := filepath.Join(tmpDir, "source.fasta")

		err := os.WriteFile(src, []byte("content"), 0644)
		assert.NoError(t, err)

		err = CopyFile(src, "/nonexistent/dir/dest.fasta")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create destination file")
	})
}
