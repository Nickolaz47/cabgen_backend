package healthservice_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/healthservice"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHealthService(t *testing.T) {
	testutils.SetupTestContext()

	country := testmodels.NewCountry("BRA", nil)
	mockHealthService := testmodels.NewHealthService(
		uuid.NewString(), "Hospitl A", models.Public, country,
		"Rio de Janeiro", "", "", "",
		false,
	)

	name, contact, email := "Hospital A", "Jão Silva", "jãos@mail.com"
	mockHealthServiceInput := models.HealthServiceUpdateInput{
		Name:         &name,
		Contactant:   &contact,
		ContactEmail: &email,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			UpdateFunc: func(
				ctx context.Context, ID uuid.UUID,
				input models.HealthServiceUpdateInput) (
				*models.HealthServiceAdminTableResponse, error) {
				return &models.HealthServiceAdminTableResponse{
					ID:           mockHealthService.ID,
					Name:         *mockHealthServiceInput.Name,
					Type:         mockHealthService.Type,
					Country:      mockHealthService.Country.Code,
					City:         mockHealthService.City,
					Contactant:   mockHealthServiceInput.Contactant,
					ContactEmail: mockHealthServiceInput.ContactEmail,
					ContactPhone: mockHealthService.ContactPhone,
					IsActive:     mockHealthService.IsActive,
				}, nil
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/health-service", testutils.ToJSON(
				mockHealthServiceInput,
			), nil, gin.Params{
				{Key: "healthServiceId", Value: mockHealthService.ID.String()},
			},
		)
		handler.UpdateHealthService(c)

		expected := testutils.ToJSON(
			map[string]any{
				"data": models.HealthServiceAdminTableResponse{
					ID:           mockHealthService.ID,
					Name:         *mockHealthServiceInput.Name,
					Type:         mockHealthService.Type,
					Country:      mockHealthService.Country.Code,
					City:         mockHealthService.City,
					Contactant:   mockHealthServiceInput.Contactant,
					ContactEmail: mockHealthServiceInput.ContactEmail,
					ContactPhone: mockHealthService.ContactPhone,
					IsActive:     mockHealthService.IsActive,
				},
			},
		)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.UpdateHealthServiceTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockHealthServiceService{}
			handler := healthservice.NewAdminHealthServiceHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPut, "/api/admin/health-service", tt.Body,
				nil, gin.Params{
					{Key: "healthServiceId", Value: uuid.NewString()},
				},
			)
			handler.UpdateHealthService(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Invalid Type", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		body := testutils.ToJSON(
			map[string]any{
				"name":         mockHealthService.Name,
				"type":         "ong",
				"country_code": mockHealthService.Country.Code,
			},
		)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/health-service", body,
			nil, gin.Params{{Key: "healthServiceId", Value: uuid.NewString()}},
		)
		handler.UpdateHealthService(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "Invalid health service type.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/health-service", "",
			nil, gin.Params{{Key: "healthServiceId", Value: "132"}},
		)
		handler.UpdateHealthService(c)

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
			UpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.HealthServiceUpdateInput) (
				*models.HealthServiceAdminTableResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/health-service", testutils.ToJSON(
				mockHealthServiceInput,
			),
			nil, gin.Params{{Key: "healthServiceId", Value: uuid.NewString()}},
		)
		handler.UpdateHealthService(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Health service not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			UpdateFunc: func(ctx context.Context, ID uuid.UUID,
				input models.HealthServiceUpdateInput) (
				*models.HealthServiceAdminTableResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/health-service", testutils.ToJSON(
				mockHealthServiceInput,
			),
			nil, gin.Params{{Key: "healthServiceId", Value: uuid.NewString()}},
		)
		handler.UpdateHealthService(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "A health service with this name already exists.",
			},
		)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			UpdateFunc: func(
				ctx context.Context, ID uuid.UUID,
				input models.HealthServiceUpdateInput) (
				*models.HealthServiceAdminTableResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		c, w := testutils.SetupGinContext(
			http.MethodPut, "/api/admin/health-service",
			testutils.ToJSON(mockHealthServiceInput), nil, gin.Params{
				{Key: "healthServiceId", Value: uuid.NewString()},
			},
		)
		handler.UpdateHealthService(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
