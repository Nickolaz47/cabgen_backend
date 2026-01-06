package sequencer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSequencer(t *testing.T) {
	testutils.SetupTestContext()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Ilumina", "MySeq", true,
	)

	brand, model := "Illumina", "MiSeq"
	mockSequencerInput := models.SequencerUpdateInput{
		Brand: &brand,
		Model: &model,
	}

	mockUpdatedSequencer := testmodels.NewSequencer(
		mockSequencer.ID.String(), brand, model, mockSequencer.IsActive,
	)

	t.Run("Success", func(t *testing.T) {
		service := testmodels.MockSequencerService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.Sequencer, error) {
				return &mockUpdatedSequencer, nil
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(mockSequencerInput),
			nil, gin.Params{{Key: "sequencerId", Value: mockSequencer.ID.String()}},
		)

		mockHandler.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockUpdatedSequencer.ToFormResponse(),
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateSequencerTests {
		t.Run(tt.Name, func(t *testing.T) {
			service := testmodels.MockSequencerService{}
			mockHandler := sequencer.NewAdminSequencerHandler(&service)

			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/sequencer", tt.Body,
				nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
			)

			mockHandler.UpdateSequencer(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Invalid ID", func(t *testing.T) {
		service := testmodels.MockSequencerService{}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: "132"}},
		)

		mockHandler.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Sequencer not found", func(t *testing.T) {
		service := testmodels.MockSequencerService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.Sequencer, error) {
				return nil, services.ErrNotFound
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(
				mockSequencerInput,
			),
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		mockHandler.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"error": "Sequencer not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("DB error", func(t *testing.T) {
		service := testmodels.MockSequencerService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.Sequencer, error) {
				return nil, services.ErrInternal
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(
				mockSequencerInput,
			),
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		mockHandler.UpdateSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
