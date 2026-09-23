package container

import (
	adminHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/audit"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildAuditService(db *gorm.DB, logger *zap.Logger) services.AuditService {
	auditRepo := repositories.NewAuditRepository(db)
	auditService := services.NewAuditService(auditRepo, logger)

	return auditService
}

func BuildAdminAuditHandler(svc services.AuditService,
) *adminHandler.AdminAuditHandler {
	return adminHandler.NewAdminAuditHandler(svc)
}
