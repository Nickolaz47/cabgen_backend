package container

import (
	adminLaboratory "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildLaboratoryService(db *gorm.DB, 
	logger *zap.Logger) services.LaboratoryService {
	labRepo := repositories.NewLaboratoryRepo(db)
	labService := services.NewLaboratoryService(labRepo, logger)

	return labService
}

func BuildLaboratoryHandler(svc services.LaboratoryService) *laboratory.LaboratoryHandler {
	return laboratory.NewLaboratoryHandler(svc)
}

func BuildAdminLaboratoryHandler(svc services.LaboratoryService) *adminLaboratory.AdminLaboratoryHandler {
	return adminLaboratory.NewAdminLaboratoryHandler(svc)
}
