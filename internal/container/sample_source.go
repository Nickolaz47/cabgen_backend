package container

import (
	adminSampleSource "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/samplesource"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildSampleSourceService(db *gorm.DB) services.SampleSourceService {
	sampleSourceRepo := repositories.NewSampleSourceRepo(db)
	sampleSourceService := services.NewSampleSourceService(sampleSourceRepo)

	return sampleSourceService
}

func BuildSampleSourceHandler(svc services.SampleSourceService) *samplesource.SampleSourceHandler {
	return samplesource.NewSampleSourceHandler(svc)
}

func BuildAdminSampleSourceHandler(svc services.SampleSourceService) *adminSampleSource.AdminSampleSourceHandler {
	return adminSampleSource.NewAdminSampleSourceHandler(svc)
}
