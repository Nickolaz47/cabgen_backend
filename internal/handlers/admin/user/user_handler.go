package user

import (
	"net/http"
	"strconv"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminUserHandler struct {
	Service services.AdminUserService
}

func NewAdminUserHandler(svc services.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{Service: svc}
}

func (h *AdminUserHandler) GetUsers(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	filter := models.AdminUserFilter{}

	if v := c.Query("input"); v != "" {
		filter.Input = &v
	}

	if v := c.Query("userRole"); v != "" {
		userRole := models.UserRole(v)
		filter.UserRole = &userRole
	}

	if v := c.Query("active"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.APIResponse{
				Error: responses.GetResponse(
					localizer, responses.AdminInvalidActiveQuery,
				),
			})
			return
		}
		filter.Active = &b
	}

	users, err := h.Service.Find(c.Request.Context(), filter, language)
	if err != nil {
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: users})
}

func (h *AdminUserHandler) GetUserByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	rawID := c.Param("userId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	user, err := h.Service.FindByID(c.Request.Context(), id, language)
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

func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var newUser models.AdminUserCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &newUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if !newUser.UserRole.IsValid() {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, responses.InvalidUserRoleError)},
		)
		return
	}

	createdUser, err := h.Service.Create(c.Request.Context(), newUser, userToken.Username, language)
	if err != nil {
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusCreated,
		responses.APIResponse{
			Data:    createdUser,
			Message: responses.GetResponse(localizer, responses.AdminRegisterSuccess),
		},
	)
}

func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	rawID := c.Param("userId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var updateInput models.AdminUserUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &updateInput); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if updateInput.UserRole != nil {
		if !updateInput.UserRole.IsValid() {
			c.JSON(http.StatusBadRequest,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.InvalidUserRoleError)},
			)
			return
		}
	}

	updateUser, err := h.Service.Update(c.Request.Context(), id, updateInput, language)
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
		Data: updateUser,
	})
}

func (h *AdminUserHandler) ActivateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("userId")

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err := h.Service.ActivateUser(c.Request.Context(), id, userToken.Username); err != nil {
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AdminUserHandler) DeactivateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("userId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err := h.Service.DeactivateUser(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleUserError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("userId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
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
			localizer, responses.UserDeleted)},
	)
}
