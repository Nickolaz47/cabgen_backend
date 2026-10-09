package container

import (
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

func BuildTemplateService(db *gorm.DB,
	countrySvc services.CountryService,
	logger *zap.Logger) services.TemplateService {
	laboratoryRepo := repositories.NewLaboratoryRepo(db)
	sequencerRepo := repositories.NewSequencerRepo(db)
	healthServiceRepo := repositories.NewHealthServiceRepo(db)
	originRepo := repositories.NewOriginRepo(db)
	microRepo := repositories.NewMicroorganismRepository(db)
	sampleSourceRepo := repositories.NewSampleSourceRepo(db)
	selectOptionsSvc := services.NewSelectOptionsService(
		laboratoryRepo, sequencerRepo, healthServiceRepo,
		originRepo, microRepo, sampleSourceRepo)

	return services.NewTemplateService(countrySvc,
		services.NewCityService(), selectOptionsSvc, logger)
}

func BuildSampleHandler(svc services.SampleService,
	templateSvc services.TemplateService) *sample.SampleHandler {
	return sample.NewSampleHandler(svc, templateSvc)
}

func BuildAdminSampleHandler(svc services.SampleService,
	templateSvc services.TemplateService) *sample.SampleHandler {
	return sample.NewAdminSampleHandler(svc, templateSvc)
}
