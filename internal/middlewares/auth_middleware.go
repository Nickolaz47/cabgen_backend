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
	tokenProvider := auth.NewTokenProvider()
	return func(c *gin.Context) {
		localizer := translation.GetLocalizerFromContext(c)

		accessSecret, err := auth.GetSecretKey(auth.Access)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
			c.Abort()
			return
		}

		tokenStr, err := auth.ExtractToken(c, auth.Access)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.UnauthorizedError)})
			c.Abort()
			return
		}

		userToken, err := tokenProvider.ValidateToken(tokenStr, accessSecret)
		if err != nil && strings.Contains(err.Error(), "token expired:") {
			c.JSON(http.StatusForbidden,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.TokenExpiredError)})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.UnauthorizedError)})
			c.Abort()
			return
		}

		c.Set("user", userToken)
		c.Next()
	}
}
