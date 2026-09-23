package laboratory

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

type AdminLaboratoryHandler struct {
	Service services.LaboratoryService
}

func NewAdminLaboratoryHandler(svc services.LaboratoryService) *AdminLaboratoryHandler {
	return &AdminLaboratoryHandler{
		Service: svc,
	}
}

func (h *AdminLaboratoryHandler) GetAllLaboratories(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesGet, nil)

	labs, err := h.Service.FindAll(c.Request.Context())
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesGetFailed, nil)
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: labs})
}

func (h *AdminLaboratoryHandler) GetLaboratoryByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesGetByID, nil)
	rawID := c.Param("laboratoryId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesGetByIDFailed,
			nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	lab, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesGetByIDFailed,
			nil)
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: lab})
}

func (h *AdminLaboratoryHandler) GetLaboratoriesByNameOrAbbreviation(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesSearch, nil)
	nameOrAbbreviation := utils.SanitizeQuery(c.Query("nameOrAbbreviation"))

	var (
		labs []models.LaboratoryAdminTableResponse
		err  error
	)

	if nameOrAbbreviation == "" {
		labs, err = h.Service.FindAll(c.Request.Context())
	} else {
		labs, err = h.Service.FindByNameOrAbbreviation(c.Request.Context(), nameOrAbbreviation)
	}

	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesSearchFailed,
			nil)
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: labs})
}

func (h *AdminLaboratoryHandler) CreateLaboratory(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesCreate, nil)
	var newLaboratory models.LaboratoryCreateInput

	if errMsg, valid := validations.Validate(c, localizer, &newLaboratory); !valid {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesCreateFailed,
			nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	lab, err := h.Service.Create(c.Request.Context(), newLaboratory)

	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesCreateFailed, nil)
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesCreate, map[string]string{"laboratory_id": lab.ID.String()})
	c.JSON(http.StatusCreated, responses.APIResponse{
		Data:    lab,
		Message: responses.GetResponse(localizer, responses.LaboratoryCreationSuccess),
	})
}

func (h *AdminLaboratoryHandler) UpdateLaboratory(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesUpdate, nil)
	rawID := c.Param("laboratoryId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesUpdateFailed,
			map[string]string{"laboratory_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var laboratoryUpdateInput models.LaboratoryUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &laboratoryUpdateInput)
	if !ok {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesUpdateFailed,
			map[string]string{"laboratory_id": rawID})
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	labUpdated, err := h.Service.Update(c.Request.Context(), id, laboratoryUpdateInput)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesUpdateFailed, map[string]string{"laboratory_id": rawID})
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesUpdate, map[string]string{"laboratory_id": rawID})

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: labUpdated,
	})
}

func (h *AdminLaboratoryHandler) DeleteLaboratory(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesDelete, nil)
	rawID := c.Param("laboratoryId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesDeleteFailed,
			map[string]string{"laboratory_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id); err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminLaboratoriesDeleteFailed,
			map[string]string{"laboratory_id": rawID})
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	validations.SetAuditEvent(c, models.AuditEventAdminLaboratoriesDelete, map[string]string{"laboratory_id": rawID})

	c.JSON(
		http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(localizer, responses.LaboratoryDeleted)},
	)
}
