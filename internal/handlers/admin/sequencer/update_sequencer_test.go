package sequencer_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
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

func TestUpdateSequencer(t *testing.T) {
	testutils.SetupTestContext()
	db := testutils.SetupTestRepos()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Ilumina", "MySeq", true,
	)
	db.Create(&mockSequencer)

	brand, model := "Illumina", "MiSeq"
	mockSequencerInput := models.SequencerUpdateInput{
		Brand: &brand,
		Model: &model,
	}

	t.Run("Success", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(mockSequencerInput),
			nil, gin.Params{{Key: "sequencerId", Value: mockSequencer.ID.String()}},
		)

		sequencer.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": models.Sequencer{
					ID:       mockSequencer.ID,
					Brand:    *mockSequencerInput.Brand,
					Model:    *mockSequencerInput.Model,
					IsActive: mockSequencer.IsActive,
				},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateSequencerTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/sequencer", tt.Body,
				nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
			)

			sequencer.UpdateSequencer(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Sequencer not found", func(t *testing.T) {
		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(
				mockSequencerInput,
			),
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		sequencer.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"error": "Sequencer not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
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
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(
				mockSequencerInput,
			),
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		sequencer.UpdateSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
