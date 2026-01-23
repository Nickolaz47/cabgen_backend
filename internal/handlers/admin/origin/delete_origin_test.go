package origin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteOrigin(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return nil
			},
		}
		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)
		handler.DeleteOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Origin deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Origin not found", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrNotFound
			},
		}
		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)
		handler.DeleteOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Origin not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockOriginService{}
		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: "123"}},
		)
		handler.DeleteOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockOriginService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return gorm.ErrInvalidTransaction
			},
		}
		handler := origin.NewAdminOriginHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)
		handler.DeleteOrigin(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
