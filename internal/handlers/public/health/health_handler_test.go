package health_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/health"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	handler := health.NewHealthHandler()

	tests := []struct {
		name     string
		status   models.PlatformStatus
		expected string
	}{
		{"ok", models.PlatformStatusOk, "ok"},
		{"platformInTesting", models.PlatformStatusPlatformInTesting, "platformInTesting"},
		{"maintenance", models.PlatformStatusMaintenance, "maintenance"},
		{"adminOnly", models.PlatformStatusAdminOnly, "adminOnly"},
		{"analysisDisabled", models.PlatformStatusAnalysisDisabled, "analysisDisabled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origStatus := config.PlatformStatus
			config.PlatformStatus = tt.status
			defer func() { config.PlatformStatus = origStatus }()

			c, w := testutils.SetupGinContext(
				http.MethodGet, "/api/health", "",
				nil, nil,
			)
			handler.Health(c)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, testutils.ToJSON(map[string]string{
				"message": tt.expected,
			}), w.Body.String())
		})
	}
}
