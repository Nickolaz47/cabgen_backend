package microorganism

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

type AdminMicroorganismHandler struct {
	Service services.MicroorganismService
}

func NewAdminMicroorganismHandler(svc services.MicroorganismService) *AdminMicroorganismHandler {
	return &AdminMicroorganismHandler{
		Service: svc,
	}
}

func (h *AdminMicroorganismHandler) GetMicroorganisms(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	micros, err := h.Service.FindAll(c.Request.Context(), language)
	if err != nil {
		code, errMsg := handlererrors.HandleMicroorganismError(err)
		c.JSON(
			code, responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Data: micros,
		})
}

func (h *AdminMicroorganismHandler) GetMicroorganismByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("microorganismId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	micro, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleMicroorganismError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: micro})
}

func (h *AdminMicroorganismHandler) GetMicroorganismBySpecies(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	species := c.Query("species")

	var (
		micros []models.MicroorganismAdminTableResponse
		err    error
	)

	if species == "" {
		micros, err = h.Service.FindAll(c.Request.Context(), language)
	} else {
		micros, err = h.Service.FindBySpecies(
			c.Request.Context(), species, language,
		)
	}

	if err != nil {
		code, errMsg := handlererrors.HandleMicroorganismError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: micros,
	})
}

func (h *AdminMicroorganismHandler) CreateMicroorganism(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newMicro models.MicroorganismCreateInput

	if errMsg, valid := validations.Validate(c, localizer, &newMicro); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if !newMicro.Taxon.IsValid() {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.MicroorganismInvalidTaxon),
			})
		return
	}

	errMsg, ok := validations.ValidateTranslationMap(
		c, "microorganism", newMicro.Variety)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	micro, err := h.Service.Create(c.Request.Context(), newMicro)
	if err != nil {
		code, errMsg := handlererrors.HandleMicroorganismError(err)
		c.JSON(code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data: micro,
		Message: responses.GetResponse(localizer,
			responses.MicroorganismCreationSuccess),
	})
}

func (h *AdminMicroorganismHandler) UpdateMicroorganism(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("microorganismId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var microUpdateInput models.MicroorganismUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &microUpdateInput)
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	if microUpdateInput.Taxon != nil {
		ok = microUpdateInput.Taxon.IsValid()
	}
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.MicroorganismInvalidTaxon),
			})
		return
	}

	if microUpdateInput.Variety != nil {
		errMsg, ok = validations.ValidateTranslationMap(
			c, "microorganism", microUpdateInput.Variety)
	}
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	microUpdated, err := h.Service.Update(c.Request.Context(), id, microUpdateInput)
	if err != nil {
		code, errMsg := handlererrors.HandleMicroorganismError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Data: microUpdated,
		})
}

func (h *AdminMicroorganismHandler) DeleteMicroorganism(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("microorganismId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleMicroorganismError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(
		http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(
			localizer, responses.MicroorganismDeleted,
		)},
	)
}
