package samplesource_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetActiveSampleSources(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockSampleSource := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Plasma", "en": "Plasma", "es": "Plasma"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		true,
	)
	mockSampleSource2 := testmodels.NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Coágulo sanguíneo", "en": "Blood clot", "es": "Coágulo de sangre"},
		map[string]string{"pt": "Sangue", "en": "Blood", "es": "Sangre"},
		false,
	)
	db.Create(&mockSampleSource)
	db.Create(&mockSampleSource2)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/sampleSource", "",
			nil, nil,
		)

		samplesource.GetActiveSampleSources(c)

		expected := testutils.ToJSON(map[string]any{
			"data": []models.SampleSourceFormResponse{
				mockSampleSource.ToFormResponse(c),
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		sampleSourceRepo := repository.SampleSourceRepo
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repository.SampleSourceRepo = repository.NewSampleSourceRepo(mockDB)
		defer func() {
			repository.SampleSourceRepo = sampleSourceRepo
		}()

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/sampleSource", "",
			nil, nil,
		)

		samplesource.GetActiveSampleSources(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
