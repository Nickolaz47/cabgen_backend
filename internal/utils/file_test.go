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
		assert.Contains(t, err.Error(),
			"failed to create destination file")
	})
}

func TestCleanupFiles(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tmpDir := t.TempDir()
		f1 := filepath.Join(tmpDir, "a.fastq")
		f2 := filepath.Join(tmpDir, "b.fastq")
		os.WriteFile(f1, []byte("data"), 0644)
		os.WriteFile(f2, []byte("data"), 0644)

		err := CleanupFiles([]string{f1, f2})
		assert.NoError(t, err)
		assert.NoFileExists(t, f1)
		assert.NoFileExists(t, f2)
	})

	t.Run("Empty List", func(t *testing.T) {
		err := CleanupFiles([]string{})
		assert.NoError(t, err)
	})

	t.Run("Missing File Returns Error", func(t *testing.T) {
		err := CleanupFiles([]string{"/nonexistent/file.fastq"})
		assert.Error(t, err)
	})

	t.Run("Partial Failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		exists := filepath.Join(tmpDir, "a.fastq")
		os.WriteFile(exists, []byte("data"), 0644)

		err := CleanupFiles([]string{exists, "/nonexistent/b.fastq"})
		assert.Error(t, err)
		assert.NoFileExists(t, exists)
	})
}

func TestResolveSampleFilePath(t *testing.T) {
	uid := "user1"
	sid := "sample1"

	t.Run("Non-FASTA Returns Direct Path", func(t *testing.T) {
		root := t.TempDir()
		path, ok := ResolveSampleFilePath(root, uid, sid,
			"reads.fastq", "fastq", "")
		assert.True(t, ok)
		assert.Equal(t,
			filepath.Join(root, "uploads", "users", uid,
				"samples", sid, "reads.fastq"),
			path)
	})

	t.Run("FASTA Found in Sample Dir", func(t *testing.T) {
		root := t.TempDir()
		sampleDir := filepath.Join(root, "uploads", "users", uid,
			"samples", sid)
		os.MkdirAll(sampleDir, 0o755)
		os.WriteFile(filepath.Join(sampleDir, "contigs.fasta"),
			[]byte(">seq\nATCG"), 0o644)

		path, ok := ResolveSampleFilePath(root, uid, sid,
			"contigs.fasta", "fasta", "analysis1")
		assert.True(t, ok)
		assert.Equal(t,
			filepath.Join(sampleDir, "contigs.fasta"), path)
	})

	t.Run("FASTA Falls Back to Assembly", func(t *testing.T) {
		root := t.TempDir()
		assemblyDir := filepath.Join(root, "uploads", "users", uid,
			"samples", sid, "analyses", "analysis1", "assembly")
		os.MkdirAll(assemblyDir, 0o755)
		os.WriteFile(filepath.Join(assemblyDir, "contigs.fasta"),
			[]byte(">seq\nATCG"), 0o644)

		path, ok := ResolveSampleFilePath(root, uid, sid,
			"contigs.fasta", "fasta", "analysis1")
		assert.True(t, ok)
		assert.Equal(t,
			filepath.Join(assemblyDir, "contigs.fasta"), path)
	})

	t.Run("FASTA Not Found Anywhere", func(t *testing.T) {
		root := t.TempDir()
		path, ok := ResolveSampleFilePath(root, uid, sid,
			"contigs.fasta", "fasta", "analysis1")
		assert.False(t, ok)
		assert.Empty(t, path)
	})

	t.Run("FASTA No Analysis ID Skips Assembly", func(t *testing.T) {
		root := t.TempDir()
		path, ok := ResolveSampleFilePath(root, uid, sid,
			"contigs.fasta", "fasta", "")
		assert.False(t, ok)
		assert.Empty(t, path)
	})
}

func TestValidateDBPath(t *testing.T) {
	tmp := t.TempDir()
	nonEmptyDir := filepath.Join(tmp, "non-empty")

	assert.NoError(t, os.Mkdir(nonEmptyDir, 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(nonEmptyDir, "db.txt"),
		[]byte("data"), 0o644))

	emptyFile := filepath.Join(tmp, "empty.txt")
	assert.NoError(t, os.WriteFile(emptyFile, nil, 0o644))

	fullFile := filepath.Join(tmp, "full.txt")
	assert.NoError(t, os.WriteFile(fullFile, []byte("refs"), 0o644))

	emptyDir := filepath.Join(tmp, "empty-dir")
	assert.NoError(t, os.Mkdir(emptyDir, 0o755))

	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{name: "Not configured", path: "", wantErr: "not configured"},
		{name: "Not found", path: filepath.Join(tmp, "missing"),
			wantErr: "not accessible"},
		{name: "Empty dir", path: emptyDir, wantErr: "is empty"},
		{name: "Non-empty dir", path: nonEmptyDir, wantErr: ""},
		{name: "Empty file", path: emptyFile, wantErr: "is empty"},
		{name: "Non-empty file", path: fullFile, wantErr: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDBPath(tt.path)

			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}
