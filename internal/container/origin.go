package container

import (
	adminOrigin "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildOriginService(db *gorm.DB) services.OriginService {
	originRepo := repository.NewOriginRepo(db)
	originService := services.NewOriginService(originRepo)

	return originService
}

func BuildOriginHandler(svc services.OriginService) *origin.OriginHandler {
	return origin.NewOriginHandler(svc)
}

func BuildAdminOriginHandler(svc services.OriginService) *adminOrigin.AdminOriginHandler {
	return adminOrigin.NewAdminOriginHandler(svc)
}
