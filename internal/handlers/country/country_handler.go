package country

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCountries(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var countries []models.Country
	if err := db.DB.Find(&countries).Error; err != nil {
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

	var country models.Country
	err := db.DB.Where("code = ?", code).First(&country).Error

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
