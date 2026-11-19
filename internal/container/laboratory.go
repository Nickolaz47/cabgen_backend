package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	adminLaboratory "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/laboratory"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
)

func BuildLaboratoryService(db db.Database,
	logger *zap.Logger) services.LaboratoryService {
	labRepo := repositories.NewLaboratoryRepo(db)
	labService := services.NewLaboratoryService(labRepo, logger)

	return labService
}

func BuildAdminLaboratoryHandler(svc services.LaboratoryService) *adminLaboratory.AdminLaboratoryHandler {
	return adminLaboratory.NewAdminLaboratoryHandler(svc)
}
