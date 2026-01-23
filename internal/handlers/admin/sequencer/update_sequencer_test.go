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
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSequencer(t *testing.T) {
	testutils.SetupTestContext()

	mockSequencer := testmodels.NewSequencer(
		uuid.NewString(), "Ilumina", "MySeq", false,
	)

	brand, model, isActive := "Illumina", "MiSeq", true
	mockSequencerInput := models.SequencerUpdateInput{
		Brand:    &brand,
		Model:    &model,
		IsActive: &isActive,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSequencerService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error) {
				return &models.SequencerAdminTableResponse{
					ID:       mockSequencer.ID,
					Model:    *mockSequencerInput.Model,
					Brand:    *mockSequencerInput.Brand,
					IsActive: *mockSequencerInput.IsActive,
				}, nil
			},
		}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(mockSequencerInput),
			nil, gin.Params{{Key: "sequencerId", Value: mockSequencer.ID.String()}},
		)
		handler.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": models.SequencerAdminTableResponse{
					ID:       mockSequencer.ID,
					Model:    *mockSequencerInput.Model,
					Brand:    *mockSequencerInput.Brand,
					IsActive: *mockSequencerInput.IsActive,
				},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateSequencerTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockSequencerService{}
			handler := sequencer.NewAdminSequencerHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/sequencer", tt.Body,
				nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
			)
			handler.UpdateSequencer(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockSequencerService{}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", "",
			nil, gin.Params{{Key: "sequencerId", Value: "132"}},
		)
		handler.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockSequencerService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(
				mockSequencerInput,
			),
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)
		handler.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"error": "Sequencer not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := &mocks.MockSequencerService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(
				mockSequencerInput,
			),
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)
		handler.UpdateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"error": "A sequencer with this model already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockSequencerService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.SequencerUpdateInput) (*models.SequencerAdminTableResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := sequencer.NewAdminSequencerHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/sequencer", testutils.ToJSON(
				mockSequencerInput,
			),
			nil, gin.Params{{Key: "sequencerId", Value: uuid.NewString()}},
		)
		handler.UpdateSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
