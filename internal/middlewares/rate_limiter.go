package middlewares

import (
	"net/http"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func GlobalRateLimitPerMinute(max int) gin.HandlerFunc {
	limiter := rate.NewLimiter(
		rate.Every(time.Minute/time.Duration(max)), max)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			localizer := translation.GetLocalizerFromContext(c)
			c.JSON(http.StatusTooManyRequests, responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.TooManyRequestsError),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
