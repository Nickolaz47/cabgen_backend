package sequencer_test

import (
	"context"
	"encoding/json"
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

	mockSequencerInput := models.SequencerCreateInput{
		Brand:    "Illumina",
		Model:    "MiSeq",
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		service := testmodels.MockSequencerService{
			CreateFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return nil
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		body := testutils.ToJSON(mockSequencerInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sequencer", body,
			nil, nil,
		)

		mockHandler.CreateSequencer(c)

		expected := testutils.ToJSON(
			map[string]any{
				"message": "Sequencer registered successfully.",
				"data": map[string]any{
					"brand":     mockSequencerInput.Brand,
					"model":     mockSequencerInput.Model,
					"is_active": mockSequencerInput.IsActive,
				},
			},
		)

		var result map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)

		if data, ok := result["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, testutils.ToJSON(result))
	})

	for _, tt := range data.CreateSequencerTests {
		t.Run(tt.Name, func(t *testing.T) {
			service := testmodels.MockSequencerService{}
			mockHandler := sequencer.NewAdminSequencerHandler(&service)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/sequencer", tt.Body,
				nil, nil,
			)

			mockHandler.CreateSequencer(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("DB error", func(t *testing.T) {
		service := testmodels.MockSequencerService{
			CreateFunc: func(ctx context.Context, sequencer *models.Sequencer) error {
				return services.ErrInternal
			},
		}
		mockHandler := sequencer.NewAdminSequencerHandler(&service)

		body := testutils.ToJSON(mockSequencerInput)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/sequencer", body,
			nil, nil,
		)

		mockHandler.CreateSequencer(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
