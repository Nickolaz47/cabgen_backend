package laboratory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateLaboratory(t *testing.T) {
	testutils.SetupTestContext()
	name, abbreviation, isActive := "Laboratório Bittar", "LB", true

	updateInput := models.LaboratoryUpdateInput{
		Name:         &name,
		Abbreviation: &abbreviation,
		IsActive:     &isActive,
	}

	mockLab := models.Laboratory{
		Name:         name,
		Abbreviation: abbreviation,
		IsActive:     isActive,
	}

	t.Run("Success", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error) {
				return &mockLab, nil
			},
		}

		handler := laboratory.NewLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/laboratory", testutils.ToJSON(updateInput),
			nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.UpdateLaboratory(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": mockLab.ToResponse(),
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateLaboratoryTests {
		t.Run(tt.Name, func(t *testing.T) {
			labSvc := MockLaboratoryService{}

			handler := laboratory.NewLaboratoryHandler(&labSvc)
			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/laboratory", tt.Body,
				nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
			)
			handler.UpdateLaboratory(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid ID", func(t *testing.T) {
		labSvc := MockLaboratoryService{}

		handler := laboratory.NewLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/laboratory", testutils.ToJSON(updateInput),
			nil, gin.Params{{Key: "laboratoryId", Value: "asdae2"}},
		)
		handler.UpdateLaboratory(c)

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
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := laboratory.NewLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/laboratory", testutils.ToJSON(updateInput),
			nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.UpdateLaboratory(c)

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
			UpdateFunc: func(ctx context.Context, ID uuid.UUID, input models.LaboratoryUpdateInput) (*models.Laboratory, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}

		handler := laboratory.NewLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/laboratory", testutils.ToJSON(updateInput),
			nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.UpdateLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
