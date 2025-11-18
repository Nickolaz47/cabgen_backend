package laboratory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteLaboratory(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return nil
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/laboratory", "",
			nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.DeleteLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"message": "Laboratory deleted successfully.",
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		labSvc := MockLaboratoryService{}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/laboratory", "",
			nil, gin.Params{{Key: "laboratoryId", Value: "asdae2"}},
		)
		handler.DeleteLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrNotFound
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/laboratory", "",
			nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.DeleteLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Laboratory not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server Error", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return gorm.ErrInvalidTransaction
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/laboratory", "",
			nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.DeleteLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
