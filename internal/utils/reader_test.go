package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func writeMockFile(t *testing.T, filePath string, data []byte) {
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Error("failed to write mock file")
	}
}
func TestLoadJSONFile(t *testing.T) {
	tempDir := t.TempDir()
	mockFile := filepath.Join(tempDir, "file.json")
	mockContent := `[{"code": "ABW", "pt": "Aruba", "en": "Aruba", "es": "Aruba"}]`

	mockErrFile := filepath.Join(tempDir, "err.json")
	mockErrContent := `[{"code": "ABW", "pt": "Aruba", "en": "Aruba", "es": "Aruba"},]`

	writeMockFile(t, mockFile, []byte(mockContent))
	writeMockFile(t, mockErrFile, []byte(mockErrContent))

	t.Run("Success", func(t *testing.T) {
		expected := []models.Country{
			{Code: "ABW", Pt: "Aruba", En: "Aruba", Es: "Aruba"},
		}
		result, err := utils.LoadJSONFile[models.Country](mockFile)

		assert.NoError(t, err)
		assert.Equal(t, expected, result, "expected structs to be equal")
	})

	t.Run("Error - File no exists", func(t *testing.T) {
		_, err := utils.LoadJSONFile[models.Country](filepath.Join(tempDir, "file2.json"))

		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot read file")
	})

	t.Run("Error - Unmarshal", func(t *testing.T) {
		_, err := utils.LoadJSONFile[models.Country](mockErrFile)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid JSON in")
	})
}
