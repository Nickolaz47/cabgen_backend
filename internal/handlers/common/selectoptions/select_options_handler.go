package selectoptions

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type SelectOptionsHandler struct {
	Service services.SelectOptionsService
}

func NewSelectOptionsHandler(
	svc services.SelectOptionsService) *SelectOptionsHandler {
	return &SelectOptionsHandler{Service: svc}
}

func (h *SelectOptionsHandler) GetSelectOptions(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	opts, err := h.Service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(
				localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: opts,
	})
}
