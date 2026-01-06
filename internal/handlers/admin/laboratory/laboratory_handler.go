package laboratory

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
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

	labs, err := h.Service.FindAll(c.Request.Context())
	if err != nil {
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
	rawID := c.Param("laboratoryId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	lab, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
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
	nameOrAbbreaviation := c.Query("nameOrAbbreaviation")

	var (
		labs []models.Laboratory
		err  error
	)

	if nameOrAbbreaviation == "" {
		labs, err = h.Service.FindAll(c.Request.Context())
	} else {
		labs, err = h.Service.FindByNameOrAbbreviation(c.Request.Context(), nameOrAbbreaviation)
	}

	if err != nil {
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
	var newLaboratory models.LaboratoryCreateInput

	if errMsg, valid := validations.Validate(c, localizer, &newLaboratory); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	labToCreate := models.Laboratory{
		Name:         newLaboratory.Name,
		Abbreviation: newLaboratory.Abbreviation,
		IsActive:     newLaboratory.IsActive,
	}

	if err := h.Service.Create(c.Request.Context(), &labToCreate); err != nil {
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data:    labToCreate.ToResponse(),
		Message: responses.GetResponse(localizer, responses.LaboratoryCreationSuccess),
	})
}

func (h *AdminLaboratoryHandler) UpdateLaboratory(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("laboratoryId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var laboratoryUpdateInput models.LaboratoryUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &laboratoryUpdateInput)
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	labUpdated, err := h.Service.Update(c.Request.Context(), id, laboratoryUpdateInput)
	if err != nil {
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: labUpdated.ToResponse(),
	})
}

func (h *AdminLaboratoryHandler) DeleteLaboratory(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("laboratoryId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleLaboratoryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(
		http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(localizer, responses.LaboratoryDeleted)},
	)
}
