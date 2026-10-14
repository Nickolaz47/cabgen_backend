package metrics

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	Service services.MetricsService
}

func NewMetricsHandler(svc services.MetricsService) *MetricsHandler {
	return &MetricsHandler{Service: svc}
}

func (h *MetricsHandler) GetUserMetrics(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UnauthorizedError),
		})
		return
	}

	metrics, err := h.Service.GetUserMetrics(c.Request.Context(),
		userToken.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: metrics,
	})
}
