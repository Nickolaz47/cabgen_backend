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
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
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

	lab := testmodels.NewLaboratory(
		uuid.NewString(),
		name,
		abbreviation,
		isActive,
	)

	adminResponse := lab.ToAdminTableResponse()

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockLaboratoryService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.LaboratoryUpdateInput,
			) (*models.LaboratoryAdminTableResponse, error) {
				return &adminResponse, nil
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/laboratory",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.UpdateLaboratory(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": adminResponse,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateLaboratoryTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := testmodels.MockLaboratoryService{}
			handler := laboratory.NewAdminLaboratoryHandler(&svc)

			c, w := testutils.SetupGinContext(
				http.MethodPut,
				"/api/admin/laboratory",
				tt.Body,
				nil,
				gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
			)
			handler.UpdateLaboratory(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := testmodels.MockLaboratoryService{}
		handler := laboratory.NewAdminLaboratoryHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/laboratory",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "laboratoryId", Value: "asdae2"}},
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
		svc := testmodels.MockLaboratoryService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.LaboratoryUpdateInput,
			) (*models.LaboratoryAdminTableResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/laboratory",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
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

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := testmodels.MockLaboratoryService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.LaboratoryUpdateInput,
			) (*models.LaboratoryAdminTableResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/laboratory",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.UpdateLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "A laboratory with this name already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := testmodels.MockLaboratoryService{
			UpdateFunc: func(
				ctx context.Context,
				ID uuid.UUID,
				input models.LaboratoryUpdateInput,
			) (*models.LaboratoryAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut,
			"/api/admin/laboratory",
			testutils.ToJSON(updateInput),
			nil,
			gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
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
