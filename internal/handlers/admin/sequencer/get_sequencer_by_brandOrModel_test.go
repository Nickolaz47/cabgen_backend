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

func TestGetSequencerByBrandOrModel(t *testing.T) {
	testutils.SetupTestContext()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Illumina", "MiSeq", true,
	)

	t.Run("Success - Brand", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			FindByBrandOrModelFunc: func(ctx context.Context, input string) ([]models.Sequencer, error) {
				return []models.Sequencer{mockSequencer}, nil
			},
		}
		handler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=illumina", "", nil, nil,
		)

		handler.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.Sequencer{mockSequencer},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Model", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			FindByBrandOrModelFunc: func(ctx context.Context, input string) ([]models.Sequencer, error) {
				return []models.Sequencer{mockSequencer}, nil
			},
		}
		handler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=miseq", "", nil, nil,
		)

		handler.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.Sequencer{mockSequencer},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Input empty", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			FindAllFunc: func(ctx context.Context) ([]models.Sequencer, error) {
				return []models.Sequencer{mockSequencer}, nil
			},
		}
		handler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=", "", nil, nil,
		)

		handler.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": []models.Sequencer{mockSequencer},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("DB error", func(t *testing.T) {
		sequencerSvc := MockSequencerService{
			FindByBrandOrModelFunc: func(ctx context.Context, input string) ([]models.Sequencer, error) {
				return nil, services.ErrInternal
			},
		}
		handler := sequencer.NewAdminSequencerHandler(&sequencerSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer/search?brandOrModel=miseq", "",
			nil, nil,
		)

		handler.GetSequencersByBrandOrModel(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
