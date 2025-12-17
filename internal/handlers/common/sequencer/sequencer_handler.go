package sequencer

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type SequencerHandler struct {
	Service services.SequencerService
}

func NewSequencerHandler(svc services.SequencerService) *SequencerHandler {
	return &SequencerHandler{
		Service: svc,
	}
}

func (h *SequencerHandler) GetActiveSequencers(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	sequencers, err := h.Service.FindAllActive(c.Request.Context())
	if err != nil {
		code, errMsg := handleError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: sequencers,
	})
}
