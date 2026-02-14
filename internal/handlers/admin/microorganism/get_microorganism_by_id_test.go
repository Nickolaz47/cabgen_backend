package microorganism_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetMicroorganismByID(t *testing.T) {
	testutils.SetupTestContext()

	mockMicro := testmodels.NewMicroorganism(
		uuid.NewString(),
		models.Bacteria,
		"Salmonella",
		map[string]string{"pt": "ssp", "en": "ssp", "es": "ssp"},
		true,
	)

	mockResponse := mockMicro.ToAdminDetailResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.MicroorganismAdminDetailResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism",
			"",
			nil,
			gin.Params{{Key: "microorganismId", Value: mockMicro.ID.String()}},
		)

		handler.GetMicroorganismByID(c)

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
			http.MethodGet,
			"/api/admin/microorganism",
			"",
			nil,
			nil,
		)

		handler.GetMicroorganismByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.MicroorganismAdminDetailResponse, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism",
			"",
			nil,
			gin.Params{{Key: "microorganismId", Value: uuid.NewString()}},
		)

		handler.GetMicroorganismByID(c)

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
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.MicroorganismAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/microorganism",
			"",
			nil,
			gin.Params{{Key: "microorganismId", Value: mockMicro.ID.String()}},
		)

		handler.GetMicroorganismByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}