package middlewares

import (
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware assigns a correlation ID to every request, returns it
// in the X-Request-ID response header and injects it into the request
// context so all downstream log lines carry it.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()

		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Request = c.Request.WithContext(
			logging.WithRequestID(c.Request.Context(), requestID))

		c.Next()
	}
}
