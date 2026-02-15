package container

import (
	adminMicroorganism "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildMicroorganismService(db *gorm.DB,
	logger *zap.Logger) services.MicroorganismService {
	microRepo := repositories.NewMicroorganismRepository(db)
	microService := services.NewMicroorganismService(microRepo, logger)

	return microService
}

func BuildMicroorganismHandler(
	svc services.MicroorganismService) *microorganism.MicroorganismHandler {
	return microorganism.NewMicroorganismHandler(svc)
}

func BuildAdminMicroorganismHandler(
	svc services.MicroorganismService) *adminMicroorganism.AdminMicroorganismHandler {
	return adminMicroorganism.NewAdminMicroorganismHandler(svc)
}
