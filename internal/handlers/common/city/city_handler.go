package city

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

type CityHandler struct {
	Service services.CityService
}

func NewCityHandler(svc services.CityService) *CityHandler {
	return &CityHandler{Service: svc}
}

func (h *CityHandler) GetCities(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	cities, err := h.Service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(
				localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: cities,
	})
}
