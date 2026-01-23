package sequencer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetSequencerByID(t *testing.T) {
	testutils.SetupTestContext()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Illumina", "MiSeq", true,
	)

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSequencerService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SequencerAdminTableResponse, error) {
				response := mockSequencer.ToAdminTableResponse()
				return &response, nil
			},
		}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: mockSequencer.ID.String()}},
		)
		handler.GetSequencerByID(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockSequencer,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockSequencerService{}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: "132"}},
		)
		handler.GetSequencerByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Sequencer not found", func(t *testing.T) {
		svc := &mocks.MockSequencerService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SequencerAdminTableResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		handler.GetSequencerByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Sequencer not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("DB error", func(t *testing.T) {
		svc := &mocks.MockSequencerService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.SequencerAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)

		handler.GetSequencerByID(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
