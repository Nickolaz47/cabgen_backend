package user

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

func GetOwnUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
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

func UpdateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var updateUser models.UpdateUserInput
	if errMsg, valid := validations.ValidateUpdateInput(c, localizer, &updateUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	var user models.User
	if err := db.DB.Preload("Country").First(&user, "id = ?", userToken.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UserNotFoundError)})
		return
	}

	validations.ApplyUpdateToUser(&user, &updateUser)

	if updateUser.CountryCode != nil {
		country, valid := validations.ValidateCountryCode(*updateUser.CountryCode)
		if !valid {
			c.JSON(http.StatusBadRequest, responses.APIResponse{
				Error: responses.GetResponse(localizer, responses.CountryNotFoundError),
			})
			return
		}
		user.Country = *country
	}

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UpdateUserError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: user.ToResponse(c),
	})
}
