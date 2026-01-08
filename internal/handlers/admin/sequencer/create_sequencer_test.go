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
	"github.com/stretchr/testify/assert"
)

func TestCreateSequencer(t *testing.T) {
	testutils.SetupTestContext()

	input := models.SequencerCreateInput{
		Brand:    "Illumina",
		Model:    "MiSeq",
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockSequencerService{
			CreateFunc: func(ctx context.Context, input models.SequencerCreateInput) (*models.SequencerAdminTableResponse, error) {
				return &models.SequencerAdminTableResponse{
					Model:    input.Model,
					Brand:    input.Brand,
					IsActive: input.IsActive,
				}, nil
			},
		}
		handler := sequencer.NewAdminSequencerHandler(&svc)

		body := testutils.ToJSON(input)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sequencer", body,
			nil, nil,
		)
		handler.CreateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"message": "Sequencer registered successfully.",
				"data": models.SequencerAdminTableResponse{
					Model:    input.Model,
					Brand:    input.Brand,
					IsActive: input.IsActive,
				},
			},
		)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.CreateSequencerTests {
		t.Run(tt.Name, func(t *testing.T) {
			service := testmodels.MockSequencerService{}
			handler := sequencer.NewAdminSequencerHandler(&service)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/sequencer", tt.Body,
				nil, nil,
			)

			handler.CreateSequencer(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := testmodels.MockSequencerService{
			CreateFunc: func(ctx context.Context, input models.SequencerCreateInput) (*models.SequencerAdminTableResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := sequencer.NewAdminSequencerHandler(&svc)

		body := testutils.ToJSON(input)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sequencer", body,
			nil, nil,
		)
		handler.CreateSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "A sequencer with this model already exists.",
		})

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := testmodels.MockSequencerService{
			CreateFunc: func(ctx context.Context, input models.SequencerCreateInput) (*models.SequencerAdminTableResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := sequencer.NewAdminSequencerHandler(&svc)

		body := testutils.ToJSON(input)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sequencer", body,
			nil, nil,
		)
		handler.CreateSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
