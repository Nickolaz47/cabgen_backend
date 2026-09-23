package audit

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type AdminAuditHandler struct {
	Service services.AuditService
}

func NewAdminAuditHandler(svc services.AuditService) *AdminAuditHandler {
	return &AdminAuditHandler{
		Service: svc,
	}
}

func (h *AdminAuditHandler) GetAuditLogs(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	var filter models.AuditFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.InvalidQueryParamError),
		})
		return
	}

	auditLogs, err := h.Service.FindAll(c.Request.Context(), filter)
	if err != nil {
		code, errMsg := handlererrors.HandleAuditError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: auditLogs})
}
