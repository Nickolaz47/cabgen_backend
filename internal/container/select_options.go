package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/selectoptions"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func BuildSelectOptionHandler(
	laboratoryRepo repositories.LaboratoryRepository,
	sequencerRepo repositories.SequencerRepository,
	healthServiceRepo repositories.HealthServiceRepository,
	originRepo repositories.OriginRepository,
	microorganismRepo repositories.MicroorganismRepository,
	sampleSourceRepo repositories.SampleSourceRepository,
) *selectoptions.SelectOptionsHandler {
	service := services.NewSelectOptionsService(
		laboratoryRepo,
		sequencerRepo,
		healthServiceRepo,
		originRepo,
		microorganismRepo,
		sampleSourceRepo,
	)
	return selectoptions.NewSelectOptionsHandler(service)
}
