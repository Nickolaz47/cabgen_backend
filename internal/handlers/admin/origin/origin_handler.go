package origin

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminOriginHandler struct {
	Service services.OriginService
}

func NewAdminOriginHandler(svc services.OriginService) *AdminOriginHandler {
	return &AdminOriginHandler{Service: svc}
}

func (h *AdminOriginHandler) GetAllOrigins(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminOriginsGet, nil)
	language := translation.GetLanguageFromContext(c)

	origins, err := h.Service.FindAll(c.Request.Context(), language)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminOriginsGetFailed, nil)
		code, errMsg := handlererrors.HandleOriginError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: origins})
}

func (h *AdminOriginHandler) GetOriginByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminOriginsGetByID, nil)
	rawID := c.Param("originId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminOriginsGetByIDFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	origin, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminOriginsGetByIDFailed, nil)
		code, errMsg := handlererrors.HandleOriginError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: origin})
}

func (h *AdminOriginHandler) GetOriginsByName(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminOriginsSearch, nil)
	language := translation.GetLanguageFromContext(c)
	name := utils.SanitizeQuery(c.Query("name"))

	var (
		origins []models.OriginAdminTableResponse
		err     error
	)

	if name == "" {
		origins, err = h.Service.FindAll(c.Request.Context(), language)
	} else {
		origins, err = h.Service.FindByName(c.Request.Context(), name, language)
	}

	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminOriginsSearchFailed, nil)
		code, errMsg := handlererrors.HandleOriginError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: origins,
	})
}

func (h *AdminOriginHandler) CreateOrigin(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminOriginsCreate, nil)

	var newOrigin models.OriginCreateInput

	if errMsg, valid := validations.Validate(c, localizer, &newOrigin); !valid {
		validations.SetAuditEvent(c, models.AuditEventAdminOriginsCreateFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	errMsg, ok := validations.ValidateTranslationMap(c, "origin", newOrigin.Names)
	if !ok {
		validations.SetAuditEvent(c, models.AuditEventAdminOriginsCreateFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	origin, err := h.Service.Create(c.Request.Context(), newOrigin)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminOriginsCreateFailed, nil)
		code, errMsg := handlererrors.HandleOriginError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data:    origin,
		Message: responses.GetResponse(localizer, responses.OriginCreationSuccess),
	})
}

func (h *AdminOriginHandler) UpdateOrigin(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminOriginsUpdate, nil)
	rawID := c.Param("originId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminOriginsUpdateFailed,
			map[string]string{"origin_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var originUpdateInput models.OriginUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &originUpdateInput)
	if !ok {
		validations.SetAuditEvent(c,
			models.AuditEventAdminOriginsUpdateFailed,
			map[string]string{"origin_id": rawID})
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	if originUpdateInput.Names != nil {
		errMsg, ok = validations.ValidateTranslationMap(c, "origin", originUpdateInput.Names)
	}

	if !ok {
		validations.SetAuditEvent(c,
			models.AuditEventAdminOriginsUpdateFailed,
			map[string]string{"origin_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	originUpdated, err := h.Service.Update(c.Request.Context(), id, originUpdateInput)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminOriginsUpdateFailed,
			map[string]string{"origin_id": rawID})
		code, errMsg := handlererrors.HandleOriginError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Data: originUpdated,
		},
	)
}

func (h *AdminOriginHandler) DeleteOrigin(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminOriginsDelete, nil)
	rawID := c.Param("originId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminOriginsDeleteFailed,
			map[string]string{"origin_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id); err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminOriginsDeleteFailed,
			map[string]string{"origin_id": rawID})
		code, errMsg := handlererrors.HandleOriginError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Message: responses.GetResponse(localizer, responses.OriginDeleted),
		})
}
