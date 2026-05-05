package healthservice_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/healthservice"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHealthService(t *testing.T) {
	testutils.SetupTestContext()

	country := testmodels.NewCountry("BRA", nil)
	mockHealthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John Doe", "john@example.com", "123456789",
		true,
	)

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return nil
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/health-service",
			"", nil, gin.Params{
				{Key: "healthServiceId", Value: mockHealthService.ID.String()},
			},
		)
		handler.DeleteHealthService(c)

		expected := testutils.ToJSON(map[string]string{
			"message": "Health service deleted successfully.",
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete, "/api/admin/health-service", "",
			nil, gin.Params{{Key: "healthServiceId", Value: "132"}},
		)
		handler.DeleteHealthService(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrNotFound
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/health-service",
			"", nil, gin.Params{
				{Key: "healthServiceId", Value: mockHealthService.ID.String()},
			},
		)
		handler.DeleteHealthService(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Health service not found.",
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			DeleteFunc: func(ctx context.Context, ID uuid.UUID) error {
				return services.ErrInternal
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodDelete,
			"/api/admin/health-service",
			"", nil, gin.Params{
				{Key: "healthServiceId", Value: mockHealthService.ID.String()},
			},
		)
		handler.DeleteHealthService(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
