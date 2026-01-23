package laboratory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllLaboratories(t *testing.T) {
	testutils.SetupTestContext()

	lab := testmodels.NewLaboratory(
		uuid.NewString(),
		"Laboratory 1",
		"LAB1",
		true,
	)

	adminResponse := lab.ToAdminTableResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockLaboratoryService{
			FindAllFunc: func(ctx context.Context) ([]models.LaboratoryAdminTableResponse, error) {
				return []models.LaboratoryAdminTableResponse{
					adminResponse,
				}, nil
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/laboratory",
			"",
			nil,
			nil,
		)
		handler.GetAllLaboratories(c)

		expected := testutils.ToJSON(
			map[string][]models.LaboratoryAdminTableResponse{
				"data": {adminResponse},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockLaboratoryService{
			FindAllFunc: func(ctx context.Context) ([]models.LaboratoryAdminTableResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := laboratory.NewAdminLaboratoryHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/laboratory",
			"",
			nil,
			nil,
		)
		handler.GetAllLaboratories(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
