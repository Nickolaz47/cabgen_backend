package sequencer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetActiveSequencers(t *testing.T) {
	testutils.SetupTestContext()
	mockSequencer := models.SequencerFormResponse{ID: uuid.New()}

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockSequencerService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.SequencerFormResponse, error) {
				return []models.SequencerFormResponse{mockSequencer}, nil
			},
		}
		handler := sequencer.NewSequencerHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/sequencer", "", nil, nil,
		)
		handler.GetActiveSequencers(c)

		expected := testutils.ToJSON(
			map[string][]models.SequencerFormResponse{
				"data": {mockSequencer},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := testmodels.MockSequencerService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.SequencerFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := sequencer.NewSequencerHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/laboratory", "",
			nil, nil,
		)
		handler.GetActiveSequencers(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
