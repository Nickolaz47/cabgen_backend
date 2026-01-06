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
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateOrigin(t *testing.T) {
	testutils.SetupTestContext()
	mockOrigin := testmodels.NewOrigin(uuid.NewString(), map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"}, true)

	t.Run("Success", func(t *testing.T) {
		originSvc := testmodels.MockOriginService{
			CreateFunc: func(ctx context.Context, origin *models.Origin) error {
				return nil
			},
		}
		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/origin", testutils.ToJSON(mockOrigin),
			nil, nil,
		)

		handler.CreateOrigin(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "Origin created successfully.",
			"data":    mockOrigin.ToResponse("en"),
		})

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.CreateOriginTests {
		t.Run(tt.Name, func(t *testing.T) {
			originSvc := testmodels.MockOriginService{}
			handler := origin.NewAdminOriginHandler(&originSvc)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/origin", tt.Body,
				nil, nil,
			)

			handler.CreateOrigin(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Conflict", func(t *testing.T) {
		originSvc := testmodels.MockOriginService{
			CreateFunc: func(ctx context.Context, origin *models.Origin) error {
				return services.ErrConflict
			},
		}
		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/origin", testutils.ToJSON(mockOrigin),
			nil, nil,
		)
		handler.CreateOrigin(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Origin already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		originSvc := testmodels.MockOriginService{
			CreateFunc: func(ctx context.Context, origin *models.Origin) error {
				return gorm.ErrInvalidTransaction
			},
		}
		handler := origin.NewAdminOriginHandler(&originSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/origin", testutils.ToJSON(mockOrigin),
			nil, nil,
		)
		handler.CreateOrigin(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
