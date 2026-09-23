package audit_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/audit"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAuditLogs(t *testing.T) {
	testutils.SetupTestContext()

	mockAudit := testmodels.NewAudit(models.AuditEventLogin, "10.0.0.1", "{}", 200)
	mockResponse := mockAudit.ToResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockAuditService{
			FindAllFunc: func(ctx context.Context, filter models.AuditFilter) (
				[]models.AuditResponse, error) {
				return []models.AuditResponse{mockResponse}, nil
			},
		}

		handler := audit.NewAdminAuditHandler(svc)
		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/audit", "", nil, nil)
		handler.GetAuditLogs(c)

		expected := testutils.ToJSON(map[string][]models.AuditResponse{
			"data": {mockResponse},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success - With Filter", func(t *testing.T) {
		svc := &mocks.MockAuditService{
			FindAllFunc: func(ctx context.Context, filter models.AuditFilter) (
				[]models.AuditResponse, error) {
				assert.Equal(t, models.AuditEventLogin, filter.Event)
				assert.Equal(t, "10.0.0.1", filter.Source)
				assert.Equal(t, 404, filter.Status)
				assert.Equal(t,
					time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), *filter.Date)
				return []models.AuditResponse{mockResponse}, nil
			},
		}

		handler := audit.NewAdminAuditHandler(svc)
		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/audit?event=auth.login&source=10.0.0.1&status=404&date=2026-01-02",
			"", nil, nil)
		handler.GetAuditLogs(c)

		expected := testutils.ToJSON(map[string][]models.AuditResponse{
			"data": {mockResponse},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid Query Params", func(t *testing.T) {
		svc := &mocks.MockAuditService{}
		handler := audit.NewAdminAuditHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/audit?date=02/01/2026", "", nil, nil)
		handler.GetAuditLogs(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "Invalid query parameters.",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Service Error", func(t *testing.T) {
		svc := &mocks.MockAuditService{
			FindAllFunc: func(ctx context.Context, filter models.AuditFilter) (
				[]models.AuditResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := audit.NewAdminAuditHandler(svc)
		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/audit", "", nil, nil)
		handler.GetAuditLogs(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
