package middlewares

import (
	"net/http"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		localizer := translation.GetLocalizerFromContext(c)

		accessSecret, err := auth.GetSecretKey(auth.Access)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
			return
		}

		userToken, err := auth.ValidateToken(c, auth.Access, accessSecret)
		if err != nil && strings.Contains(err.Error(), "token expired:") {
			c.JSON(http.StatusForbidden,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.TokenExpiredError)})
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.UnauthorizedError)})
			return
		}

		c.Set("user", userToken)
		c.Next()
	}
}
