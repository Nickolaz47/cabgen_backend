package country

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type PublicCountryHandler struct {
	Service services.CountryService
}

func NewPublicCountryHandler(svc services.CountryService) *PublicCountryHandler {
	return &PublicCountryHandler{Service: svc}
}

func (h *PublicCountryHandler) GetCountries(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	countries, err := h.Service.FindAll(c.Request.Context(), language)
	if err != nil {
		code, errMsg := handlererrors.HandleCountryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: countries})
}
