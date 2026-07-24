package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/selectoptions"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func BuildSelectOptionHandler(
	laboratoryService services.LaboratoryService,
	sequencerService services.SequencerService,
	healthServiceService services.HealthServiceService,
	originService services.OriginService,
	microorganismService services.MicroorganismService,
	sampleSourceService services.SampleSourceService,
) *selectoptions.SelectOptionsHandler {
	service := services.NewSelectOptionsService(
		laboratoryService,
		sequencerService,
		healthServiceService,
		originService,
		microorganismService,
		sampleSourceService,
	)
	return selectoptions.NewSelectOptionsHandler(service)
}
