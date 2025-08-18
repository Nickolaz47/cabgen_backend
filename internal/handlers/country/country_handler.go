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

var CountryRepo *repository.CountryRepository

func GetCountries(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	countries, err := CountryRepo.GetCountries()
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

	country, err := CountryRepo.GetCountry(code)
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
