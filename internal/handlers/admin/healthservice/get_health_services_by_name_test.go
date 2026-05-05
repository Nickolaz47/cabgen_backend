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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthServicesByName(t *testing.T) {
	testutils.SetupTestContext()

	country := testmodels.NewCountry("BRA", nil)
	mockHealthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospital A", models.Public, country,
		"Rio de Janeiro", "John Doe", "john@example.com", "123456789",
		true,
	)
	mockHealthService2 := testmodels.NewHealthService(
		uuid.NewString(), "Laboratório B", models.Public, country,
		"Rio de Janeiro", "", "", "",
		true,
	)

	t.Run("Success - Not empty name", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			FindAllFunc: func(ctx context.Context) (
				[]models.HealthServiceAdminTableResponse, error) {
				return []models.HealthServiceAdminTableResponse{
					mockHealthService.ToAdminTableResponse(),
					mockHealthService2.ToAdminTableResponse(),
				}, nil
			},
			FindByNameFunc: func(ctx context.Context, name string) (
				[]models.HealthServiceAdminTableResponse, error) {
				return []models.HealthServiceAdminTableResponse{
					mockHealthService.ToAdminTableResponse(),
				}, nil
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/health-service/search?name=hospital",
			"", nil, nil,
		)
		handler.GetHealthServicesByName(c)

		expected := testutils.ToJSON(map[string]any{
			"data": []models.HealthServiceAdminTableResponse{
				mockHealthService.ToAdminTableResponse(),
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - Empty name", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			FindAllFunc: func(ctx context.Context) (
				[]models.HealthServiceAdminTableResponse, error) {
				return []models.HealthServiceAdminTableResponse{
					mockHealthService.ToAdminTableResponse(),
					mockHealthService2.ToAdminTableResponse(),
				}, nil
			},
			FindByNameFunc: func(ctx context.Context, name string) (
				[]models.HealthServiceAdminTableResponse, error) {
				return []models.HealthServiceAdminTableResponse{
					mockHealthService.ToAdminTableResponse(),
				}, nil
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/health-service/search?name=",
			"", nil, nil,
		)
		handler.GetHealthServicesByName(c)

		expected := testutils.ToJSON(map[string]any{
			"data": []models.HealthServiceAdminTableResponse{
				mockHealthService.ToAdminTableResponse(),
				mockHealthService2.ToAdminTableResponse(),
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			FindByNameFunc: func(ctx context.Context, name string) (
				[]models.HealthServiceAdminTableResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodGet,
			"/api/admin/health-service/search?name=hospital",
			"", nil, nil,
		)
		handler.GetHealthServicesByName(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
