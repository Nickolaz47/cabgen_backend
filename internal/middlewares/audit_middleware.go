package middlewares

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

func AuditMiddleware(auditSvc services.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(validations.AuditKey,
			&models.AuditInput{Source: c.ClientIP()})

		c.Next()

		audit, ok := validations.GetAuditInputFromContext(c)
		if !ok || audit.Event == "" {
			return
		}

		audit.Status = c.Writer.Status()
		if userToken, ok := validations.GetUserTokenFromContext(c); ok {
			audit.UserID = &userToken.ID
		}

		go func() {
			auditSvc.Create(context.WithoutCancel(c.Request.Context()), audit)
		}()
	}
}
