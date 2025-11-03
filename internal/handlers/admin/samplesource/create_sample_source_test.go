package samplesource_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateSampleSource(t *testing.T) {
	testutils.SetupTestContext()
	testutils.SetupTestRepos()

	mockSampleSourceInput := models.SampleSourceCreateInput{
		Names:    map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		Groups:   map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sampleSource", body,
			nil, nil,
		)

		samplesource.CreateSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "Sample source created successfully.",
			"data": map[string]any{
				"name":      mockSampleSourceInput.Names["en"],
				"group":     mockSampleSourceInput.Groups["en"],
				"is_active": mockSampleSourceInput.IsActive,
			},
		})

		var result map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)

		if data, ok := result["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, testutils.ToJSON(result))
	})

	for _, tt := range data.CreateSampleSourceTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/sampleSource", tt.Body,
				nil, nil,
			)

			samplesource.CreateSampleSource(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("DB error", func(t *testing.T) {
		origRepo := repository.SampleSourceRepo
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repository.SampleSourceRepo = repository.NewSampleSourceRepo(mockDB)
		defer func() {
			repository.SampleSourceRepo = origRepo
		}()

		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sampleSource", body,
			nil, nil,
		)

		samplesource.CreateSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
