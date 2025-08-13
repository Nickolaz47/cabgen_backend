package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(consoleLogger, fileLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logMsg := fmt.Sprintf(
			"status=%d method=%s path=%s client_ip=%s latency=%s",
			status,
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			latency,
		)

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("latency", latency),
		}

		if gin.Mode() == gin.DebugMode {
			switch {
			case status >= 500:
				consoleLogger.Error("Server Error - " + logMsg)
			case status >= 400:
				consoleLogger.Warn("Client Error - " + logMsg)
			default:
				consoleLogger.Info("Request processed - " + logMsg)
			}
		}

		switch {
		case status >= 500:
			fileLogger.Error("Server Error", fields...)
		case status >= 400:
			fileLogger.Warn("Client Error", fields...)
		default:
			fileLogger.Info("Request processed", fields...)
		}
	}
}
