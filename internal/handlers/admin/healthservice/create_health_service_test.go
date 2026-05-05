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
	"github.com/stretchr/testify/assert"
)

func TestCreateHealthService(t *testing.T) {
	testutils.SetupTestContext()

	input := models.HealthServiceCreateInput{
		Name:        "Laboratório Dom",
		Type:        models.Private,
		CountryCode: "BRA",
		IsActive:    true,
	}

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			CreateFunc: func(
				ctx context.Context, input models.HealthServiceCreateInput) (
				*models.HealthServiceAdminTableResponse, error) {
				return &models.HealthServiceAdminTableResponse{
					Name:     input.Name,
					Type:     input.Type,
					Country:  input.CountryCode,
					IsActive: input.IsActive,
				}, nil
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		body := testutils.ToJSON(input)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/health-service", body,
			nil, nil,
		)
		handler.CreateHealthService(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "Health service created successfully.",
			"data": models.HealthServiceAdminTableResponse{
				Name:     input.Name,
				Type:     input.Type,
				Country:  input.CountryCode,
				IsActive: input.IsActive,
			},
		})

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Type", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		body := testutils.ToJSON(
			map[string]any{
				"name":         input.Name,
				"type":         "ong",
				"country_code": input.CountryCode,
			},
		)

		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/health-service", body,
			nil, nil,
		)
		handler.CreateHealthService(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "Invalid health service type.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.CreateHealthServiceTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockHealthServiceService{}
			handler := healthservice.NewAdminHealthServiceHandler(svc)

			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/admin/health-service", tt.Body,
				nil, nil,
			)

			handler.CreateHealthService(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Conflict", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			CreateFunc: func(
				ctx context.Context, input models.HealthServiceCreateInput) (
				*models.HealthServiceAdminTableResponse, error) {
				return nil, services.ErrConflict
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		body := testutils.ToJSON(input)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/health-service", body,
			nil, nil,
		)
		handler.CreateHealthService(c)

		expected := testutils.ToJSON(map[string]any{
			"error": "A health service with this name already exists.",
		})

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal", func(t *testing.T) {
		svc := &mocks.MockHealthServiceService{
			CreateFunc: func(
				ctx context.Context, input models.HealthServiceCreateInput) (
				*models.HealthServiceAdminTableResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := healthservice.NewAdminHealthServiceHandler(svc)

		body := testutils.ToJSON(input)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/admin/health-service", body,
			nil, nil,
		)
		handler.CreateHealthService(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
