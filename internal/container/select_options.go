package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/selectoptions"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func BuildSelectOptionHandler() *selectoptions.SelectOptionsHandler {
	service := services.NewSelectOptionsService()
	return selectoptions.NewSelectOptionsHandler(service)
}
