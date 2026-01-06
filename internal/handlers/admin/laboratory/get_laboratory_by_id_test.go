package laboratory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetLaboratoryByID(t *testing.T) {
	testutils.SetupTestContext()
	mockLab := testmodels.NewLaboratory(
		uuid.NewString(), "Laboratory 1", "LAB1", true,
	)

	t.Run("Success", func(t *testing.T) {
		labSvc := testmodels.MockLaboratoryService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return &mockLab, nil
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/laboratory", "",
			nil, gin.Params{{Key: "laboratoryId", Value: mockLab.ID.String()}},
		)
		handler.GetLaboratoryByID(c)

		expected := testutils.ToJSON(
			map[string]models.Laboratory{
				"data": mockLab,
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		labSvc := testmodels.MockLaboratoryService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, nil
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/laboratory", "",
			nil, nil,
		)
		handler.GetLaboratoryByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		labSvc := testmodels.MockLaboratoryService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, services.ErrNotFound
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/laboratory", "",
			nil, gin.Params{{Key: "laboratoryId", Value: uuid.NewString()}},
		)
		handler.GetLaboratoryByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Laboratory not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server Error", func(t *testing.T) {
		labSvc := testmodels.MockLaboratoryService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
				return nil, services.ErrInternal
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/admin/laboratory", "",
			nil, gin.Params{{Key: "laboratoryId", Value: mockLab.ID.String()}},
		)
		handler.GetLaboratoryByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
