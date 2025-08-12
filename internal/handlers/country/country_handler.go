package country

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCountries(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	countries, err := repository.GetCountryRepo().GetCountries()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: countries})
}

func GetCountryByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	code := c.Param("code")

	country, err := repository.GetCountryRepo().GetCountry(code)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.CountryNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: country})
}
