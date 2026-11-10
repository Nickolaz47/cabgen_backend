package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const expectedAdminOnlyMessage = "Currently, only administrators have access to this resource."

func TestPlatformAccessMiddleware(t *testing.T) {
	handlerCalled := false
	handler := func(c *gin.Context) { handlerCalled = true; c.Status(http.StatusOK) }

	tests := []struct {
		name       string
		status     models.PlatformStatus
		path       string
		method     string
		wantStatus int
		wantCalled bool
	}{
		{"ok allows any", models.PlatformStatusOk, "/api/contact", http.MethodGet, http.StatusOK, true},
		{"platformInTesting no impact", models.PlatformStatusPlatformInTesting, "/api/contact", http.MethodGet, http.StatusOK, true},
		{"adminOnly no impact on protected route", models.PlatformStatusAdminOnly, "/api/analyses", http.MethodGet, http.StatusOK, true},
		{"maintenance blocks public route", models.PlatformStatusMaintenance, "/api/contact", http.MethodGet, http.StatusServiceUnavailable, false},
		{"maintenance allows health", models.PlatformStatusMaintenance, "/api/health", http.MethodGet, http.StatusOK, true},
		{"maintenance allows metrics", models.PlatformStatusMaintenance, "/api/metrics", http.MethodGet, http.StatusOK, true},
		{"maintenance allows countries", models.PlatformStatusMaintenance, "/api/countries", http.MethodGet, http.StatusOK, true},
		{"maintenance allows OPTIONS", models.PlatformStatusMaintenance, "/api/contact", http.MethodOptions, http.StatusOK, true},
		{"analysisDisabled no effect on non-analysis", models.PlatformStatusAnalysisDisabled, "/api/contact", http.MethodGet, http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origStatus := config.PlatformStatus
			config.PlatformStatus = tt.status
			defer func() { config.PlatformStatus = origStatus }()
			handlerCalled = false

			c, w := testutils.SetupGinContext(tt.method, tt.path, "", nil, nil)
			middlewares.PlatformAccessMiddleware()(c)
			if !c.IsAborted() {
				handler(c)
			}

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantCalled, handlerCalled)
			if !tt.wantCalled && tt.status == models.PlatformStatusMaintenance {
				assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
			}
		})
	}
}

func TestAdminOnlyPublicBlockMiddleware(t *testing.T) {
	handlerCalled := false
	handler := func(c *gin.Context) { handlerCalled = true; c.Status(http.StatusOK) }

	tests := []struct {
		name       string
		status     models.PlatformStatus
		path       string
		wantStatus int
		wantCalled bool
	}{
		{"ok allows public route", models.PlatformStatusOk, "/api/contact", http.StatusOK, true},
		{"adminOnly blocks register", models.PlatformStatusAdminOnly, "/api/auth/register", http.StatusServiceUnavailable, false},
		{"adminOnly blocks contact", models.PlatformStatusAdminOnly, "/api/contact", http.StatusServiceUnavailable, false},
		{"adminOnly allows login", models.PlatformStatusAdminOnly, "/api/auth/login", http.StatusOK, true},
		{"adminOnly allows logout", models.PlatformStatusAdminOnly, "/api/auth/logout", http.StatusOK, true},
		{"adminOnly allows refresh", models.PlatformStatusAdminOnly, "/api/auth/refresh", http.StatusOK, true},
		{"adminOnly allows health", models.PlatformStatusAdminOnly, "/api/health", http.StatusOK, true},
		{"adminOnly allows metrics", models.PlatformStatusAdminOnly, "/api/metrics", http.StatusOK, true},
		{"adminOnly allows countries", models.PlatformStatusAdminOnly, "/api/countries", http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origStatus := config.PlatformStatus
			config.PlatformStatus = tt.status
			defer func() { config.PlatformStatus = origStatus }()
			handlerCalled = false

			c, w := testutils.SetupGinContext(http.MethodGet, tt.path, "", nil, nil)
			middlewares.AdminOnlyPublicBlockMiddleware()(c)
			if !c.IsAborted() {
				handler(c)
			}

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantCalled, handlerCalled)
			if !tt.wantCalled {
				assert.Contains(t, w.Body.String(), expectedAdminOnlyMessage)
				assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
			}
		})
	}
}

func TestAdminOnlyMiddleware(t *testing.T) {
	handlerCalled := false
	handler := func(c *gin.Context) { handlerCalled = true; c.Status(http.StatusOK) }

	tests := []struct {
		name       string
		status     models.PlatformStatus
		role       models.UserRole
		wantCalled bool
	}{
		{"ok no check", models.PlatformStatusOk, models.Collaborator, true},
		{"adminOnly allows admin", models.PlatformStatusAdminOnly, models.Admin, true},
		{"adminOnly blocks collaborator", models.PlatformStatusAdminOnly, models.Collaborator, false},
		{"adminOnly blocks no token", models.PlatformStatusAdminOnly, "", false},
		{"adminOnly blocks collaborator on me", models.PlatformStatusAdminOnly, models.Collaborator, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origStatus := config.PlatformStatus
			config.PlatformStatus = tt.status
			defer func() { config.PlatformStatus = origStatus }()
			handlerCalled = false

			path := "/api/samples"
			if tt.name == "adminOnly blocks collaborator on me" {
				path = "/api/auth/me"
			}
			c, w := testutils.SetupGinContext(http.MethodGet, path, "", nil, nil)
			if tt.role != "" {
				c.Set(validations.UserTokenKey, &models.UserToken{UserRole: tt.role})
			}

			middlewares.AdminOnlyMiddleware()(c)
			if !c.IsAborted() {
				handler(c)
			}

			if tt.wantCalled {
				assert.True(t, handlerCalled)
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.False(t, handlerCalled)
				assert.Equal(t, http.StatusServiceUnavailable, w.Code)
				assert.Contains(t, w.Body.String(), expectedAdminOnlyMessage)
				assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
			}
		})
	}
}

func TestAnalysisDisabledMiddleware(t *testing.T) {
	handlerCalled := false
	handler := func(c *gin.Context) { handlerCalled = true; c.Status(http.StatusOK) }

	tests := []struct {
		name       string
		status     models.PlatformStatus
		method     string
		wantCalled bool
	}{
		{"ok allows POST", models.PlatformStatusOk, http.MethodPost, true},
		{"analysisDisabled allows GET", models.PlatformStatusAnalysisDisabled, http.MethodGet, true},
		{"analysisDisabled blocks POST", models.PlatformStatusAnalysisDisabled, http.MethodPost, false},
		{"analysisDisabled blocks DELETE", models.PlatformStatusAnalysisDisabled, http.MethodDelete, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origStatus := config.PlatformStatus
			config.PlatformStatus = tt.status
			defer func() { config.PlatformStatus = origStatus }()
			handlerCalled = false

			c, w := testutils.SetupGinContext(tt.method, "/api/analyses", "", nil, nil)
			middlewares.AnalysisDisabledMiddleware()(c)
			if !c.IsAborted() {
				handler(c)
			}

			if tt.wantCalled {
				assert.True(t, handlerCalled)
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.False(t, handlerCalled)
				assert.Equal(t, http.StatusServiceUnavailable, w.Code)
			}
		})
	}
}
