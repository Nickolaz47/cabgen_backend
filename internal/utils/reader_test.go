package utils_test

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoadJSONFile(t *testing.T) {
	tempDir := t.TempDir()
	mockFile := filepath.Join(tempDir, "file.json")
	mockContent := `[{"code": "ABW", "Names": {"pt": "Aruba", "en": "Aruba", "es": "Aruba"}}]`

	mockErrFile := filepath.Join(tempDir, "err.json")
	mockErrContent := `[{"code": "ABW", "Names": {"pt": "Aruba", "en": "Aruba", "es": "Aruba"}},]`

	testutils.WriteMockFile(t, mockFile, []byte(mockContent))
	testutils.WriteMockFile(t, mockErrFile, []byte(mockErrContent))

	t.Run("Success", func(t *testing.T) {
		expected := []models.Country{
			{
				Code:  "ABW",
				Names: map[string]string{
					"pt": "Aruba", "en": "Aruba", "es": "Aruba",
				}},
		}
		result, err := utils.LoadJSONFile[models.Country](mockFile)

		assert.NoError(t, err)
		assert.Equal(t, expected, result,
			"expected structs to be equal")
	})

	t.Run("Error - File no exists", func(t *testing.T) {
		_, err := utils.LoadJSONFile[models.Country](
			filepath.Join(tempDir, "file2.json"))

		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot read file")
	})

	t.Run("Error - Unmarshal", func(t *testing.T) {
		_, err := utils.LoadJSONFile[models.Country](mockErrFile)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid JSON in")
	})
}

func TestIsGzip(t *testing.T) {
	t.Run("Valid Gzip", func(t *testing.T) {
		data := append([]byte{0x1f, 0x8b}, []byte("junk")...)
		r, ok := utils.IsGzip(bytes.NewReader(data))
		assert.True(t, ok)

		result, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("Invalid Gzip", func(t *testing.T) {
		data := []byte("not gzip")
		r, ok := utils.IsGzip(bytes.NewReader(data))
		assert.False(t, ok)

		result, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("Too Short", func(t *testing.T) {
		data := []byte{0x1f}
		r, ok := utils.IsGzip(bytes.NewReader(data))
		assert.False(t, ok)

		result, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x1f, 0x00}, result)
	})

	t.Run("Empty", func(t *testing.T) {
		r, ok := utils.IsGzip(bytes.NewReader(nil))
		assert.False(t, ok)

		result, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x00}, result)
	})

	t.Run("Partial Read Then Invalid", func(t *testing.T) {
		data := []byte("xxff")
		r, ok := utils.IsGzip(strings.NewReader("xxff"))
		assert.False(t, ok)

		result, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})
}
