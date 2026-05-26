package container

import (
	adminSample "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildSampleService(db *gorm.DB, rootDir string,
	logger *zap.Logger) services.SampleService {
	sampleRepo := repositories.NewSampleRepo(db)
	countryRepo := repositories.NewCountryRepo(db)
	userRepo := repositories.NewUserRepo(db)
	originRepo := repositories.NewOriginRepo(db)
	sampleSourceRepo := repositories.NewSampleSourceRepo(db)
	microRepo := repositories.NewMicroorganismRepository(db)
	sequencerRepo := repositories.NewSequencerRepo(db)
	labRepo := repositories.NewLaboratoryRepo(db)
	healthServiceRepo := repositories.NewHealthServiceRepo(db)

	sampleService := services.NewSampleService(
		sampleRepo, countryRepo, userRepo, originRepo,
		sampleSourceRepo, microRepo, sequencerRepo, labRepo,
		healthServiceRepo, rootDir, logger,
	)

	return sampleService
}

func BuildSampleHandler(svc services.SampleService) *sample.SampleHandler {
	return sample.NewSampleHandler(svc)
}

func BuildAdminSampleHandler(
	svc services.SampleService) *adminSample.AdminSampleHandler {
	return adminSample.NewAdminSampleHandler(svc)
}
