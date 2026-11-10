package middlewares

import (
	"net/http"
	"slices"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

var (
	maintenancePublicPaths = []string{
		"/api/health",
		"/api/metrics",
		"/api/countries",
	}
	adminOnlyPublicPaths = append(maintenancePublicPaths,
		"/api/auth/login", "/api/auth/logout", "/api/auth/refresh")
)

func PlatformAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Let CORS preflight requests through regardless of platform status.
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		if config.PlatformStatus == models.PlatformStatusMaintenance &&
			!containsPath(maintenancePublicPaths, path) {
			clearAuthCookies(c)
			c.JSON(http.StatusServiceUnavailable,
				responses.APIResponse{Message: models.PlatformStatusMaintenance.String()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func adminOnlyMessage(c *gin.Context) string {
	localizer := translation.GetLocalizerFromContext(c)
	return responses.GetResponse(localizer, responses.AdminOnlyMessage)
}

func AdminOnlyPublicBlockMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.PlatformStatus != models.PlatformStatusAdminOnly {
			c.Next()
			return
		}

		if containsPath(adminOnlyPublicPaths, c.Request.URL.Path) {
			c.Next()
			return
		}

		clearAuthCookies(c)
		c.JSON(http.StatusServiceUnavailable,
			responses.APIResponse{Error: adminOnlyMessage(c)})
		c.Abort()
	}
}

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.PlatformStatus != models.PlatformStatusAdminOnly {
			c.Next()
			return
		}

		userToken, ok := validations.GetUserTokenFromContext(c)
		if !ok || userToken.UserRole != models.Admin {
			clearAuthCookies(c)
			c.JSON(http.StatusServiceUnavailable,
				responses.APIResponse{Error: adminOnlyMessage(c)})
			c.Abort()
			return
		}

		c.Next()
	}
}

func clearAuthCookies(c *gin.Context) {
	http.SetCookie(c.Writer, auth.DeleteCookie(auth.Access, "/"))
	http.SetCookie(c.Writer, auth.DeleteCookie(auth.Refresh,
		"/api/auth/refresh"))
}

func AnalysisDisabledMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.PlatformStatus != models.PlatformStatusAnalysisDisabled {
			c.Next()
			return
		}

		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusServiceUnavailable,
				responses.APIResponse{
					Message: models.PlatformStatusAnalysisDisabled.String()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func containsPath(paths []string, path string) bool {
	if slices.Contains(paths, path) {
		return true
	}

	return false
}
