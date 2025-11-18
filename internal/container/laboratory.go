package container

import (
	adminLaboratory "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildLaboratoryService(db *gorm.DB) services.LaboratoryService {
	labRepo := repository.NewLaboratoryRepo(db)
	labService := services.NewLaboratoryService(labRepo)

	return labService
}

func BuildLaboratoryHandler(svc services.LaboratoryService) *laboratory.LaboratoryHandler {
	return laboratory.NewLaboratoryHandler(svc)
}

func BuildAdminLaboratoryHandler(svc services.LaboratoryService) *adminLaboratory.AdminLaboratoryHandler {
	return adminLaboratory.NewAdminLaboratoryHandler(svc)
}
