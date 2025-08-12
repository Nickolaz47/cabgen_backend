package user

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	user, err := repository.GetUserRepo().GetUserByID(userToken.ID)
	if err != nil {
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
	if errMsg, valid := validations.Validate(c, localizer, &updateUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	user, err := repository.GetUserRepo().GetUserByID(userToken.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UserNotFoundError)})
		return
	}

	if updateUser.Username != nil {
		_, err = repository.GetUserRepo().GetUserByUsername(*updateUser.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
			)
			return
		}

		if err == nil {
			c.JSON(http.StatusConflict,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterUsernameAlreadyExistsError)},
			)
			return
		}
	}

	validations.ApplyUpdateToUser(user, &updateUser)

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

	if err := repository.GetUserRepo().UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UpdateUserError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: user.ToResponse(c),
	})
}
