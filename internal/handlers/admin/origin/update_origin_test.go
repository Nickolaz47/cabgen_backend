package origin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateOrigin(t *testing.T) {
	testutils.SetupTestContext()
	names, isActive := map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true

	updateInput := models.OriginUpdateInput{
		Names:    names,
		IsActive: &isActive,
	}

	mockOrigin := models.Origin{
		Names:    names,
		IsActive: isActive,
	}

	t.Run("Success", func(t *testing.T) {
		originSvc := MockOriginService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error) {
				return &mockOrigin, nil
			},
		}
		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", testutils.ToJSON(updateInput),
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)
		handler.UpdateOrigin(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockOrigin.ToResponse(""),
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateOriginTests {
		t.Run(tt.Name, func(t *testing.T) {
			originSvc := MockOriginService{}
			handler := origin.NewAdminOriginHandler(&originSvc)

			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/origin", tt.Body,
				nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
			)
			handler.UpdateOrigin(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid ID", func(t *testing.T) {
		originSvc := MockOriginService{}
		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", "",
			nil, gin.Params{{Key: "originId", Value: "12"}},
		)
		handler.UpdateOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		originSvc := MockOriginService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", testutils.ToJSON(updateInput),
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)
		handler.UpdateOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{"error": "Origin not found."},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		originSvc := MockOriginService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.OriginUpdateInput) (*models.Origin, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/origin", testutils.ToJSON(mockOrigin),
			nil, gin.Params{{Key: "originId", Value: uuid.NewString()}},
		)
		handler.UpdateOrigin(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
