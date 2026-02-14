package microorganism_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteMicroorganism(t *testing.T) {
	testutils.SetupTestContext()

	validID := uuid.NewString()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return nil
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/microorganism",
			"",
			nil,
			gin.Params{{Key: "microorganismId", Value: validID}},
		)

		handler.DeleteMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]any{
				"message": "Microorganism deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{}
		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/microorganism",
			"",
			nil,
			nil,
		)

		handler.DeleteMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockMicroorganismService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrNotFound
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/microorganism",
			"",
			nil,
			gin.Params{{Key: "microorganismId", Value: validID}},
		)

		handler.DeleteMicroorganism(c)

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
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return gorm.ErrInvalidTransaction
			},
		}

		handler := microorganism.NewAdminMicroorganismHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/microorganism",
			"",
			nil,
			gin.Params{{Key: "microorganismId", Value: validID}},
		)

		handler.DeleteMicroorganism(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
