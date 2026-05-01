package healthservice_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/healthservice"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetActiveHealthServices(t *testing.T) {
	testutils.SetupTestContext()
	mockHealthService := models.HealthServiceFormResponse{
		ID: uuid.New(), Name: "Hospital Universitário",
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			FindAllActiveFunc: func(
				ctx context.Context) ([]models.HealthServiceFormResponse, error) {
				return []models.HealthServiceFormResponse{mockHealthService}, nil
			},
		}
		handler := healthservice.NewHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/health-service", "",
			nil, nil,
		)
		handler.GetActiveHealthServices(c)

		expected := testutils.ToJSON(
			map[string][]models.HealthServiceFormResponse{
				"data": {mockHealthService},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			FindAllActiveFunc: func(
				ctx context.Context) ([]models.HealthServiceFormResponse, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		handler := healthservice.NewHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet, "/api/health-service", "",
			nil, nil,
		)
		handler.GetActiveHealthServices(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
