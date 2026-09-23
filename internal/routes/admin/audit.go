package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/audit"
	"github.com/gin-gonic/gin"
)

func SetupAdminAuditRoutes(r *gin.RouterGroup,
	handler *audit.AdminAuditHandler) {
	auditRouter := r.Group("/audit")

	auditRouter.GET("", handler.GetAuditLogs)
}
