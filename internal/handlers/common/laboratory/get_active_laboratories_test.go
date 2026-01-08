package laboratory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetActiveLaboratories(t *testing.T) {
	testutils.SetupTestContext()
	mockLab := models.LaboratoryFormResponse{ID: uuid.New()}

	t.Run("Success", func(t *testing.T) {
		svc := testmodels.MockLaboratoryService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
				return []models.LaboratoryFormResponse{mockLab}, nil
			},
		}
		handler := laboratory.NewLaboratoryHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/laboratory", "",
			nil, nil,
		)
		handler.GetActiveLaboratories(c)

		expected := testutils.ToJSON(
			map[string][]models.LaboratoryFormResponse{
				"data": {mockLab},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := testmodels.MockLaboratoryService{
			FindAllActiveFunc: func(ctx context.Context) ([]models.LaboratoryFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := laboratory.NewLaboratoryHandler(&svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/laboratory", "",
			nil, nil,
		)
		handler.GetActiveLaboratories(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
