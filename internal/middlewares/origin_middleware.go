package middlewares

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func OriginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		origin := c.GetHeader("Origin")
		if origin != "" && origin != config.FrontendURL {
			localizer := translation.GetLocalizerFromContext(c)
			c.JSON(http.StatusForbidden, responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.ForbiddenError),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
