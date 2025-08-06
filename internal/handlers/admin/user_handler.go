package admin

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

func GetAllUsers(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var users []models.User
	if err := db.DB.Preload("Country").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: users})
}

func GetUserByUsername(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	username := c.Param("username")

	var user models.User
	err := db.DB.Preload("Country").Where("username = ?", username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.UserNotFoundError)})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: user})
}
