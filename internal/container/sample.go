package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
)

func BuildSampleService(db db.Database, rootDir string,
	logger *zap.Logger) services.SampleService {
	sampleRepo := repositories.NewSampleRepo(db.DB())
	countryRepo := repositories.NewCountryRepo(db.DB())
	userRepo := repositories.NewUserRepo(db.DB())
	originRepo := repositories.NewOriginRepo(db.DB())
	sampleSourceRepo := repositories.NewSampleSourceRepo(db.DB())
	microRepo := repositories.NewMicroorganismRepository(db.DB())
	sequencerRepo := repositories.NewSequencerRepo(db.DB())
	labRepo := repositories.NewLaboratoryRepo(db)
	healthServiceRepo := repositories.NewHealthServiceRepo(db.DB())

	sampleService := services.NewSampleService(
		sampleRepo, countryRepo, userRepo, originRepo,
		sampleSourceRepo, microRepo, sequencerRepo, labRepo,
		healthServiceRepo, rootDir, logger,
	)

	return sampleService
}

func BuildTemplateService(db db.Database,
	countrySvc services.CountryService,
	logger *zap.Logger) services.TemplateService {
	laboratoryRepo := repositories.NewLaboratoryRepo(db)
	sequencerRepo := repositories.NewSequencerRepo(db.DB())
	healthServiceRepo := repositories.NewHealthServiceRepo(db.DB())
	originRepo := repositories.NewOriginRepo(db.DB())
	microRepo := repositories.NewMicroorganismRepository(db.DB())
	sampleSourceRepo := repositories.NewSampleSourceRepo(db.DB())
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
