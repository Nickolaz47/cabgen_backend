package auth

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
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

func (h *AuthHandler) Me(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	meResponse := models.MeResponse{
		ID:       userToken.ID,
		Username: userToken.Username,
		UserRole: userToken.UserRole,
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: meResponse,
	})
}
