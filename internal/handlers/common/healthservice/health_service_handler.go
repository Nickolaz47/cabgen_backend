package healthservice

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type HealthServiceHandler struct {
	Service services.HealthServiceService
}

func NewHealthServiceHandler(
	svc services.HealthServiceService) *HealthServiceHandler {
	return &HealthServiceHandler{
		Service: svc,
	}
}

func (h *HealthServiceHandler) GetActiveHealthServices(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	healthServices, err := h.Service.FindAllActive(c.Request.Context())
	if err != nil {
		code, errMsg := handlererrors.HandleHealthServiceError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: healthServices,
	})
}
