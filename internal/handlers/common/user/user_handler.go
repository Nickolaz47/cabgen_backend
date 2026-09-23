package user

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
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
	validations.SetAuditEvent(c, models.AuditEventUsersGet, nil)
	language := translation.GetLanguageFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		validations.SetAuditEvent(c, models.AuditEventUsersGetFailed, nil)
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	user, err := h.Service.FindByID(c.Request.Context(), userToken.ID, language)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventUsersGetFailed, nil)
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
	validations.SetAuditEvent(c, models.AuditEventUsersUpdate, nil)
	language := translation.GetLanguageFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		validations.SetAuditEvent(c, models.AuditEventUsersUpdateFailed, nil)
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var updateUser models.UserUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &updateUser); !valid {
		validations.SetAuditEvent(c, models.AuditEventUsersUpdateFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	updatedUser, err := h.Service.Update(c.Request.Context(), userToken.ID,
		updateUser, language)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventUsersUpdateFailed, nil)
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

func (h *UserHandler) DeleteUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventUsersDelete, nil)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		validations.SetAuditEvent(c, models.AuditEventUsersDeleteFailed, nil)
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	if err := h.Service.Delete(c.Request.Context(), userToken.ID); err != nil {
		validations.SetAuditEvent(c, models.AuditEventUsersDeleteFailed, nil)
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	accessCookie := auth.DeleteCookie(auth.Access, "/")
	refreshCookie := auth.DeleteCookie(auth.Refresh, "/api/auth/refresh")
	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(
			localizer, responses.UserSelfDeleted)},
	)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventUsersUpdatePassword, nil)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		validations.SetAuditEvent(c, models.AuditEventUsersUpdatePasswordFailed, nil)
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var input models.UpdatePasswordInput
	if errMsg, valid := validations.Validate(c, localizer, &input); !valid {
		validations.SetAuditEvent(c, models.AuditEventUsersUpdatePasswordFailed,
			nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if err := h.Service.UpdatePassword(c.Request.Context(),
		userToken.ID, input); err != nil {
		validations.SetAuditEvent(c, models.AuditEventUsersUpdatePasswordFailed,
			nil)
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(
			localizer, responses.UpdatePasswordSuccess)},
	)
}

func (h *UserHandler) RequestEmailUpdate(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventUsersRequestEmailUpdate, nil)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		validations.SetAuditEvent(c,
			models.AuditEventUsersRequestEmailUpdateFailed,
			nil)
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var input models.RequestEmailUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &input); !valid {
		validations.SetAuditEvent(c,
			models.AuditEventUsersRequestEmailUpdateFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if err := h.Service.RequestEmailUpdate(c.Request.Context(),
		userToken.ID, input); err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventUsersRequestEmailUpdateFailed, nil)
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(
			localizer, responses.RequestEmailUpdateSuccess)},
	)
}

func (h *UserHandler) ConfirmEmailUpdate(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventUsersConfirmEmailUpdate, nil)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		validations.SetAuditEvent(c,
			models.AuditEventUsersConfirmEmailUpdateFailed,
			nil)
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var input models.ConfirmEmailUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &input); !valid {
		validations.SetAuditEvent(c,
			models.AuditEventUsersConfirmEmailUpdateFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if err := h.Service.ConfirmEmailUpdate(c.Request.Context(),
		userToken.ID, input); err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventUsersConfirmEmailUpdateFailed, nil)
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(
			localizer, responses.ConfirmEmailUpdateSuccess)},
	)
}
