package middlewares

import (
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

		var logAt func(*zap.Logger, string, ...zap.Field)
		var msg string
		switch {
		case status >= 500:
			logAt, msg = (*zap.Logger).Error, "Server Error"
		case status >= 400:
			logAt, msg = (*zap.Logger).Warn, "Client Error"
		default:
			logAt, msg = (*zap.Logger).Info, "Request processed"
		}

		fields := []zap.Field{
			zap.String("request_id", c.GetString("request_id")),
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("latency", latency),
		}

		if gin.Mode() == gin.DebugMode {
			logAt(consoleLogger, msg, fields...)
		}
		logAt(fileLogger, msg, fields...)
	}
}
