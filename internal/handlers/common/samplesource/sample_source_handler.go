package samplesource

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type SampleSourceHandler struct {
	Service services.SampleSourceService
}

func NewSampleSourceHandler(svc services.SampleSourceService) *SampleSourceHandler {
	return &SampleSourceHandler{Service: svc}
}

func (h *SampleSourceHandler) GetActiveSampleSources(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	activeSamplesSources, err := h.Service.FindAllActive(c.Request.Context(), language)
	if err != nil {
		code, errMsg := handlererrors.HandleOriginError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: activeSamplesSources,
	})
}
