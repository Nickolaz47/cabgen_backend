package microorganism_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateMicroorganism(t *testing.T) {
	testutils.SetupTestContext()

	mockMicro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Virus,
		"Influenza A",
		map[string]string{"pt": "H1N1", "en": "H1N1", "es": "H1N1"},
		true,
	)

	mockResponse := mockMicro.ToAdminDetailResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.MicroorganismUpdateInput) (*models.MicroorganismAdminDetailResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		input := map[string]any{
			"taxon":     "Virus",
			"species":   "Influenza A",
			"variety":   map[string]string{"pt": "H1N1", "en": "H1N1", "es": "H1N1"},
			"is_active": true,
		}

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/microorganism",
			testutils.ToJSON(input),
			nil,
			gin.Params{{Key: "microorganismId", Value: mockMicro.ID.String()}},
		)

		handler.UpdateMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockResponse,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{}
		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/microorganism",
			"",
			nil,
			nil,
		)

		handler.UpdateMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Bad Request", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{}
		handler := microorganism.NewAdminMicroorganismHandler(svc)

		for _, test := range data.UpdateMicroorganismTests {
			t.Run(test.Name, func(t *testing.T) {
				c, w := testutils.SetupGinContext(
					http.MethodPut,
					"/api/admin/microorganism",
					test.Body,
					nil,
					gin.Params{{Key: "microorganismId", Value: mockMicro.ID.String()}},
				)

				handler.UpdateMicroorganism(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, test.Expected, w.Body.String())
			})
		}
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.MicroorganismUpdateInput) (*models.MicroorganismAdminDetailResponse, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		input := map[string]any{
			"taxon": "Virus",
		}

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/microorganism",
			testutils.ToJSON(input),
			nil,
			gin.Params{{Key: "microorganismId", Value: uuid.NewString()}},
		)

		handler.UpdateMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Microorganism not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.MicroorganismUpdateInput) (*models.MicroorganismAdminDetailResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		input := map[string]any{
			"taxon": "Virus",
		}

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/microorganism",
			testutils.ToJSON(input),
			nil,
			gin.Params{{Key: "microorganismId", Value: mockMicro.ID.String()}},
		)

		handler.UpdateMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
