package samplesource_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdateSampleSource(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasm", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sange"},
		false,
	)
	db.Create(&mockSampleSource)

	isActive := true
	mockSampleSourceInput := models.SampleSourceUpdateInput{
		Names:    map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		Groups:   map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sampleSource", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
		)

		samplesource.UpdateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": map[string]any{
					"name":      mockSampleSourceInput.Names["en"],
					"group":     mockSampleSourceInput.Groups["en"],
					"is_active": *mockSampleSourceInput.IsActive,
				},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateSampleSourceTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/sampleSource", tt.Body,
				nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
			)

			samplesource.UpdateSampleSource(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Invalid ID", func(t *testing.T) {
		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sampleSource", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: "123"}},
		)

		samplesource.UpdateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Sample source not found", func(t *testing.T) {
		body := testutils.ToJSON(mockSampleSourceInput)
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sampleSource", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)

		samplesource.UpdateSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sample source not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

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
			http.MethodPut, "/api/admin/sampleSource", body,
			nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
		)

		samplesource.CreateSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
