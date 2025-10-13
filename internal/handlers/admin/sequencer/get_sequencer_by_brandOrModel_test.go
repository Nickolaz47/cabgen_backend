package sequencer_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetSequencerByBrandOrModel(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Illumina", "MiSeq", true,
	)
	db.Create(&mockSequencer)

	t.Run("Success - Brand", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=illumina", "", nil, nil,
		)

		sequencer.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.Sequencer{mockSequencer},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Model", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=miseq", "", nil, nil,
		)

		sequencer.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.Sequencer{mockSequencer},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Query parameter empty", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=", "", nil, nil,
		)

		sequencer.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(
			map[string]any{
				"error": "The search parameter brandOrModel is empty.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("DB error", func(t *testing.T) {
		origRepo := repository.SequencerRepo
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repository.SequencerRepo = repository.NewSequencerRepo(mockDB)
		defer func() {
			repository.SequencerRepo = origRepo
		}()

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=miseq", "",
			nil, nil,
		)

		sequencer.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
