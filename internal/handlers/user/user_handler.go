package user

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func GetOwnUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	rawUserToken, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	userToken, ok := rawUserToken.(*models.UserToken)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var user models.User
	if err := db.DB.Preload("Country").Where("id = ?", userToken.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: user.ToResponse(c)})
}

func UpdateUser(c *gin.Context) {}
