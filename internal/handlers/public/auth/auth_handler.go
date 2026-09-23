package auth

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

type AuthHandler struct {
	Service services.AuthService
}

func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{Service: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventRegister, nil)

	var newUser models.UserRegisterInput
	if errMsg, valid := validations.Validate(c, localizer, &newUser); !valid {
		validations.SetAuditEvent(c, models.AuditEventRegisterFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	response, err := h.Service.Register(c.Request.Context(), newUser, language)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventRegisterFailed,
			map[string]string{"auth_identity": newUser.Email})
		code, errMsg := handlererrors.HandleAuthError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusCreated,
		responses.APIResponse{
			Data:    response,
			Message: responses.GetResponse(localizer, responses.RegisterMessage),
		},
	)
}

func (h *AuthHandler) Login(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventLogin, nil)

	var login models.LoginInput
	if errMsg, valid := validations.Validate(c, localizer, &login); !valid {
		validations.SetAuditEvent(c, models.AuditEventLoginFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	cookies, err := h.Service.Login(c.Request.Context(), login)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventLoginFailed,
			map[string]string{"auth_identity": login.Username})
		code, errMsg := handlererrors.HandleAuthError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	http.SetCookie(c.Writer, cookies.AccessCookie)
	http.SetCookie(c.Writer, cookies.RefreshCookie)

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Message: responses.GetResponse(localizer, responses.LoginSuccess)})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventRefreshToken, nil)

	tokenStr, err := auth.ExtractToken(c, auth.Refresh)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventRefreshTokenFailed, nil)
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{
				Error: responses.GetResponse(
					localizer, responses.UnauthorizedError)})
		return
	}

	accessCookie, err := h.Service.Refresh(c.Request.Context(), tokenStr)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventRefreshTokenFailed, nil)
		code, errMsg := handlererrors.HandleAuthError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	http.SetCookie(c.Writer, accessCookie)

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(localizer,
			responses.TokenRenewed)})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventLogout, nil)

	accessCookie := auth.DeleteCookie(auth.Access, "/")
	refreshCookie := auth.DeleteCookie(auth.Refresh, "/api/auth/refresh")

	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Message: responses.GetResponse(localizer, responses.LogoutSuccess)},
	)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventForgotPassword, nil)

	var input models.ForgotPasswordInput
	if errMsg, valid := validations.Validate(c, localizer, &input); !valid {
		validations.SetAuditEvent(c, models.AuditEventForgotPasswordFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if err := h.Service.ForgotPassword(c.Request.Context(), input); err != nil {
		validations.SetAuditEvent(c, models.AuditEventForgotPasswordFailed,
			map[string]string{"auth_identity": input.Email})
		code, errMsg := handlererrors.HandleAuthError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer,
			responses.ForgotPasswordSuccess),
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventResetPassword, nil)

	var input models.ResetPasswordInput
	if errMsg, valid := validations.Validate(c, localizer, &input); !valid {
		validations.SetAuditEvent(c, models.AuditEventResetPasswordFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if err := h.Service.ResetPassword(c.Request.Context(), input); err != nil {
		validations.SetAuditEvent(c, models.AuditEventResetPasswordFailed, nil)
		code, errMsg := handlererrors.HandleAuthError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer,
			responses.ResetPasswordSuccess,
		),
	})
}
