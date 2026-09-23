package audit_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/audit"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAuditSelectOptions(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAuditService{
			FindAuditSelectOptionsFunc: func(ctx context.Context) (
				*models.AuditSelectOptionsResponse, error) {
				return &models.AuditSelectOptionsResponse{
					Events: []models.SelectOption{
						{Label: "auth.login", Value: "auth.login"},
					},
					Users: []models.SelectOption{
						{Label: "admin", Value: "123e4567-e89b-12d3-a456-426614174000"},
					},
				}, nil
			},
		}

		handler := audit.NewAdminAuditHandler(svc)
		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/audit/select-options", "", nil, nil)
		handler.GetAuditSelectOptions(c)

		expected := testutils.ToJSON(map[string]any{
			"data": map[string]any{
				"events": []map[string]string{
					{"label": "auth.login", "value": "auth.login"},
				},
				"users": []map[string]string{
					{"label": "admin",
						"value": "123e4567-e89b-12d3-a456-426614174000"},
				},
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Service Error", func(t *testing.T) {
		svc := &mocks.MockAuditService{
			FindAuditSelectOptionsFunc: func(ctx context.Context) (
				*models.AuditSelectOptionsResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := audit.NewAdminAuditHandler(svc)
		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/audit/select-options", "", nil, nil)
		handler.GetAuditSelectOptions(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
