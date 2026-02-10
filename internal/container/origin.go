package container

import (
	adminOrigin "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildOriginService(db *gorm.DB, logger *zap.Logger) services.OriginService {
	originRepo := repositories.NewOriginRepo(db)
	originService := services.NewOriginService(originRepo, logger)

	return originService
}

func BuildOriginHandler(svc services.OriginService) *origin.OriginHandler {
	return origin.NewOriginHandler(svc)
}

func BuildAdminOriginHandler(svc services.OriginService) *adminOrigin.AdminOriginHandler {
	return adminOrigin.NewAdminOriginHandler(svc)
}
