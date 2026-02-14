package microorganism

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type MicroorganismHandler struct {
	Service services.MicroorganismService
}

func NewMicroorganismHandler(svc services.MicroorganismService) *MicroorganismHandler {
	return &MicroorganismHandler{
		Service: svc,
	}
}

func (h *MicroorganismHandler) GetActiveMicroorganisms(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	micros, err := h.Service.FindAllActive(c.Request.Context(), language)
	if err != nil {
		code, errMsg := handlererrors.HandleMicroorganismError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: micros,
	})
}
