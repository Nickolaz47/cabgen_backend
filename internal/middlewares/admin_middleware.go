package middlewares

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		localizer := translation.GetLocalizerFromContext(c)

		userToken, ok := validations.GetUserTokenFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized,
				responses.APIResponse{Error: responses.GetResponse(localizer,
					responses.UnauthorizedError)})
			c.Abort()
			return
		}

		if userToken.UserRole != models.Admin {
			c.JSON(http.StatusUnauthorized,
				responses.APIResponse{Error: responses.GetResponse(localizer,
					responses.UnauthorizedError)})
			c.Abort()
			return
		}

		c.Next()
	}
}
