package samplesource_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeleteSampleSource(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	db.Create(&mockSampleSource)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sampleSource", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
		)

		samplesource.DeleteSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Sample source deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sampleSource", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: "12"}},
		)

		samplesource.DeleteSampleSource(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Sample source not found", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sampleSource", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: uuid.NewString()}},
		)

		samplesource.DeleteSampleSource(c)

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

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sampleSource", "",
			nil, gin.Params{{Key: "sampleSourceId", Value: mockSampleSource.ID.String()}},
		)

		samplesource.DeleteSampleSource(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
