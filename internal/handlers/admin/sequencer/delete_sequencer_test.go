package sequencer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSequencer(t *testing.T) {
	testutils.SetupTestContext()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Illumina", "MiSeq", true,
	)

	t.Run("Success", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return nil
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: mockSequencer.ID.String()}},
		)

		mockHandler.DeleteSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "Sequencer deleted successfully.",
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Invalid ID", func(t *testing.T) {
		sequencerSvc := MockSequencerService{}
		mockHandler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: "132"}},
		)

		mockHandler.DeleteSequencer(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Sequencer not found", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrNotFound
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		mockHandler.DeleteSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "Sequencer not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("DB Error", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrInternal
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		mockHandler.DeleteSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
