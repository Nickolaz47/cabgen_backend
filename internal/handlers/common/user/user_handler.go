package user

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service services.UserService
}

func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{Service: svc}
}

func (h *UserHandler) GetOwnUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	user, err := h.Service.FindByID(c.Request.Context(), userToken.ID, language)
	if err != nil {
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var updateUser models.UserUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &updateUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	updatedUser, err := h.Service.Update(c.Request.Context(), userToken.ID,
		updateUser, language)
	if err != nil {
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: updatedUser,
	})
}
