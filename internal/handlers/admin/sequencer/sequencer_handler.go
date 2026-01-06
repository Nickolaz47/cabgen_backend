package sequencer

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

type AdminSequencerHandler struct {
	Service services.SequencerService
}

func NewAdminSequencerHandler(svc services.SequencerService) *AdminSequencerHandler {
	return &AdminSequencerHandler{Service: svc}
}

func (h *AdminSequencerHandler) GetAllSequencers(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	sequencers, err := h.Service.FindAll(c.Request.Context())
	if err != nil {
		code, errMsg := handlererrors.HandleSequencerError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sequencers})
}

func (h *AdminSequencerHandler) GetSequencerByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sequencerId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	sequencer, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleSequencerError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sequencer})
}

func (h *AdminSequencerHandler) GetSequencersByBrandOrModel(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	brandOrModel := c.Query("brandOrModel")

	var (
		sequencers []models.Sequencer
		err        error
	)

	if brandOrModel == "" {
		sequencers, err = h.Service.FindAll(c.Request.Context())
	} else {
		sequencers, err = h.Service.FindByBrandOrModel(c.Request.Context(), brandOrModel)
	}

	if err != nil {
		code, errMsg := handlererrors.HandleSequencerError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sequencers})
}

func (h *AdminSequencerHandler) CreateSequencer(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newSequencer models.SequencerCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &newSequencer); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	sequencerToCreate := models.Sequencer{
		Brand:    newSequencer.Brand,
		Model:    newSequencer.Model,
		IsActive: newSequencer.IsActive,
	}

	if err := h.Service.Create(c.Request.Context(), &sequencerToCreate); err != nil {
		code, errMsg := handlererrors.HandleSequencerError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.SequencerCreationSuccess),
		Data:    sequencerToCreate,
	})
}

func (h *AdminSequencerHandler) UpdateSequencer(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sequencerId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var sequencerUpdateInput models.SequencerUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &sequencerUpdateInput)
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	sequencerUpdated, err := h.Service.Update(c.Request.Context(), id, sequencerUpdateInput)
	if err != nil {
		code, errMsg := handlererrors.HandleSequencerError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sequencerUpdated.ToFormResponse()})
}

func (h *AdminSequencerHandler) DeleteSequencer(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sequencerId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleSequencerError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.SequencerDeleted),
	})
}
