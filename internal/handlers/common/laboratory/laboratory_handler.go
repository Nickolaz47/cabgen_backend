package laboratory

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type LaboratoryHandler struct {
	Service services.LaboratoryService
}

func NewLaboratoryHandler(svc services.LaboratoryService) *LaboratoryHandler {
	return &LaboratoryHandler{
		Service: svc,
	}
}

func (h *LaboratoryHandler) GetActiveLaboratories(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	labs, err := h.Service.FindAllActive(c.Request.Context())
	if err != nil {
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: labs,
	})
}
