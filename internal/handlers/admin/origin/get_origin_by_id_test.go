package origin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetOriginByID(t *testing.T) {
	testutils.SetupTestContext()

	mockOrigin := testmodels.NewOrigin(
		uuid.NewString(),
		map[string]string{"pt": "Alimentar", "en": "Food", "es": "Alimentaria"},
		true,
	)

	mockResponse := mockOrigin.ToAdminDetailResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.OriginAdminDetailResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin",
			"",
			nil,
			gin.Params{{Key: "originId", Value: mockOrigin.ID.String()}},
		)

		handler.GetOriginByID(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockResponse,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockOriginService{}
		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin",
			"",
			nil,
			nil,
		)

		handler.GetOriginByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.OriginAdminDetailResponse, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin",
			"",
			nil,
			gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)

		handler.GetOriginByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Origin not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.OriginAdminDetailResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/origin",
			"",
			nil,
			gin.Params{{Key: "originId", Value: mockOrigin.ID.String()}},
		)

		handler.GetOriginByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
