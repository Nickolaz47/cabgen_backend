package samplesource

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

type AdminSampleSourceHandler struct {
	Service services.SampleSourceService
}

func NewAdminSampleSourceHandler(svc services.SampleSourceService) *AdminSampleSourceHandler {
	return &AdminSampleSourceHandler{Service: svc}
}

func (h *AdminSampleSourceHandler) GetSampleSources(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	sampleSources, err := h.Service.FindAll(c.Request.Context(), language)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleSourceError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sampleSources})
}

func (h *AdminSampleSourceHandler) GetSampleSourceByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleSourceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	sampleSource, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleSourceError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sampleSource})
}

func (h *AdminSampleSourceHandler) GetSampleSourcesByNameOrGroup(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	input := c.Query("nameOrGroup")

	var (
		sampleSources []models.SampleSourceAdminTableResponse
		err           error
	)

	if input == "" {
		sampleSources, err = h.Service.FindAll(c.Request.Context(), language)
	} else {
		sampleSources, err = h.Service.FindByNameOrGroup(
			c.Request.Context(),
			input,
			language,
		)
	}

	if err != nil {
		code, errMsg := handlererrors.HandleSampleSourceError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sampleSources})
}

func (h *AdminSampleSourceHandler) CreateSampleSource(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var input models.SampleSourceCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &input); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if errMsg, ok := validations.ValidateTranslationMap(c, "sampleSource", input.Names); !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if errMsg, ok := validations.ValidateTranslationMap(c, "sampleSource", input.Groups); !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	sampleSource, err := h.Service.Create(c.Request.Context(), input)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleSourceError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data:    sampleSource,
		Message: responses.GetResponse(localizer, responses.SampleSourceCreationSuccess),
	})
}

func (h *AdminSampleSourceHandler) UpdateSampleSource(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleSourceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var input models.SampleSourceUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &input)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if input.Names != nil {
		if errMsg, ok = validations.ValidateTranslationMap(c, "sampleSource", input.Names); !ok {
			c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
			return
		}
	}

	if input.Groups != nil {
		if errMsg, ok = validations.ValidateTranslationMap(c, "sampleSource", input.Groups); !ok {
			c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
			return
		}
	}

	updated, err := h.Service.Update(c.Request.Context(), id, input)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleSourceError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: updated,
	})
}

func (h *AdminSampleSourceHandler) DeleteSampleSource(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleSourceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleSampleSourceError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.SampleSourceDeleted),
	})
}
