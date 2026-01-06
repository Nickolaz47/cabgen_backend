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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateLaboratory(t *testing.T) {
	testutils.SetupTestContext()
	mockLab := testmodels.NewLaboratory(
		uuid.NewString(), "Laboratory 1", "LAB1", true,
	)

	t.Run("Success", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			CreateFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return nil
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/laboratory", testutils.ToJSON(mockLab),
			nil, nil,
		)
		handler.CreateLaboratory(c)

		expected := testutils.ToJSON(
			map[string]any{
				"message": "Laboratory registered successfully.",
				"data":    mockLab.ToResponse(),
			},
		)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.CreateLaboratoryTests {
		t.Run(tt.Name, func(t *testing.T) {
			labSvc := MockLaboratoryService{}
			handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/laboratory", tt.Body,
				nil, nil,
			)
			handler.CreateLaboratory(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Conflict", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			CreateFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return services.ErrConflict
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/laboratory", testutils.ToJSON(mockLab),
			nil, nil,
		)
		handler.CreateLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "A laboratory with this name already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		labSvc := MockLaboratoryService{
			CreateFunc: func(ctx context.Context, lab *models.Laboratory) error {
				return gorm.ErrInvalidTransaction
			},
		}

		handler := laboratory.NewAdminLaboratoryHandler(&labSvc)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/laboratory", testutils.ToJSON(mockLab),
			nil, nil,
		)
		handler.CreateLaboratory(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
