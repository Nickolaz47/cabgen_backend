package origin

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type OriginHandler struct {
	Service services.OriginService
}

func NewOriginHandler(svc services.OriginService) *OriginHandler {
	return &OriginHandler{
		Service: svc,
	}
}

func (h *OriginHandler) GetActiveOrigins(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	origins, err := h.Service.FindAllActive(c.Request.Context(), language)
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
		Data: origins,
	})
}
