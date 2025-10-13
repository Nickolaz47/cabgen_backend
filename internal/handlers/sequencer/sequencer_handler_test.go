package sequencer_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
)

func TestGetActiveSequencers(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Illumina", "MiSeq", true,
	)
	mockSequencer2 := testmodels.NewSequencer(
		uuid.NewString(), "Nanopore", "MinION", false,
	)
	db.Create(&mockSequencer)
	db.Create(&mockSequencer2)

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(http.MethodGet, "/api/sequencer", "", nil, nil)

		sequencer.GetActiveSequencers(c)

		expected := testutils.ToJSON(map[string]any{"data": []models.Sequencer{mockSequencer}})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		origRepo := repository.SequencerRepo
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repository.SequencerRepo = repository.NewSequencerRepo(mockDB)
		defer func() {
			repository.SequencerRepo = origRepo
		}()

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer", "",
			nil, nil,
		)

		sequencer.GetActiveSequencers(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
