package microorganism_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateMicroorganism(t *testing.T) {
	testutils.SetupTestContext()

	mockMicro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Escherichia coli",
		map[string]string{"pt": "Padrão", "en": "Standard", "es": "Estándar"},
		true,
	)

	mockResponse := mockMicro.ToAdminDetailResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			CreateFunc: func(ctx context.Context, input models.MicroorganismCreateInput) (*models.MicroorganismAdminDetailResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		input := map[string]any{
			"taxon":     "Bacteria",
			"species":   "Escherichia coli",
			"variety":   map[string]string{"pt": "Padrão", "en": "Standard", "es": "Estándar"},
			"is_active": true,
		}

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/microorganism",
			testutils.ToJSON(input),
			nil,
			nil,
		)

		handler.CreateMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data":    mockResponse,
				"message": "Microorganism created successfully.",
			},
		)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{}
		handler := microorganism.NewAdminMicroorganismHandler(svc)

		for _, test := range data.CreateMicroorganismTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPost,
					"/api/admin/microorganism",
					test.Body,
					nil,
					nil,
				)

				handler.CreateMicroorganism(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			CreateFunc: func(ctx context.Context, input models.MicroorganismCreateInput) (*models.MicroorganismAdminDetailResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		input := map[string]any{
			"taxon":     "Bacteria",
			"species":   "Escherichia coli",
			"variety":   map[string]string{"pt": "Padrão", "en": "Standard", "es": "Estándar"},
			"is_active": true,
		}

		c, w := testutils.SetupGinContext(
			http.MethodPost,
			"/api/admin/microorganism",
			testutils.ToJSON(input),
			nil,
			nil,
		)

		handler.CreateMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
