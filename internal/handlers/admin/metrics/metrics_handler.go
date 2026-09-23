package metrics

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

type AdminMetricsHandler struct {
	Service services.MetricsService
}

func NewAdminMetricsHandler(svc services.MetricsService) *AdminMetricsHandler {
	return &AdminMetricsHandler{
		Service: svc,
	}
}

func (h *AdminMetricsHandler) GetMetrics(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminMetrics, nil)

	metrics, err := h.Service.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: metrics})
}
