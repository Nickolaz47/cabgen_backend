package admin

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
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