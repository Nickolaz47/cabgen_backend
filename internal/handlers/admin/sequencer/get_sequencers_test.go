package sequencer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllSequencers(t *testing.T) {
	testutils.SetupTestContext()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Illumina", "MiSeq", true,
	)

	t.Run("Success", func(t *testing.T) {
		service := testmodels.MockSequencerService{
			FindAllFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return []models.Sequencer{mockSequencer}, nil
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		c, w := testutils.SetupGinContext(http.MethodGet, "/api/admin/sequencer", "", nil, nil)
		mockHandler.GetAllSequencers(c)

		expected := testutils.ToJSON(map[string]any{"data": []models.Sequencer{mockSequencer}})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		service := testmodels.MockSequencerService{
			FindAllFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return nil, services.ErrInternal
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer", "",
			nil, nil,
		)

		mockHandler.GetAllSequencers(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
