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

	var newUser models.UserRegisterInput
	if errMsg, valid := validations.Validate(c, localizer, &newUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	response, err := h.Service.Register(c.Request.Context(), newUser, language)
	if err != nil {
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

	var login models.LoginInput
	if errMsg, valid := validations.Validate(c, localizer, &login); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	cookies, err := h.Service.Login(c.Request.Context(), login)
	if err != nil {
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

func (h *AuthHandler) Logout(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	accessCookie := auth.DeleteCookie(auth.Access, "/")
	refreshCookie := auth.DeleteCookie(auth.Refresh, "/api/auth/refresh")

	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Message: responses.GetResponse(localizer, responses.LogoutSuccess)},
	)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	tokenStr, err := auth.ExtractToken(c, auth.Refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{
				Error: responses.GetResponse(
					localizer, responses.UnauthorizedError)})
		return
	}

	accessCookie, err := h.Service.Refresh(c.Request.Context(), tokenStr)
	if err != nil {
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
		responses.APIResponse{Message: responses.GetResponse(localizer, responses.TokenRenewed)})
}
