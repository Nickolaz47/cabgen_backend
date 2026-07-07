package analysis

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

type AdminAnalysisHandler struct {
	Service services.AdminAnalysisService
}

func NewAdminAnalysisHandler(svc services.AdminAnalysisService,
) *AdminAnalysisHandler {
	return &AdminAnalysisHandler{
		Service: svc,
	}
}

func (h *AdminAnalysisHandler) GetAnalyses(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	analyses, err := h.Service.FindAll(c.Request.Context())
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: analyses})
}

func (h *AdminAnalysisHandler) GetAnalysisByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	analysis, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: analysis})
}

func (h *AdminAnalysisHandler) GetAnalysisTypes(c *gin.Context) {
	c.JSON(http.StatusOK, responses.APIResponse{
		Data: models.AnalysisTypes,
	})
}

func (h *AdminAnalysisHandler) CreateAnalysis(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newAnalysis models.AdminAnalysisCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &newAnalysis); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if !newAnalysis.Type.IsValid() {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.AnalysisInvalidType),
		})
		return
	}

	payload := models.AnalysisCreateDTO(newAnalysis)
	analysis, err := h.Service.Create(c.Request.Context(), payload)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data: analysis,
		Message: responses.GetResponse(localizer,
			responses.AnalysisCreationSuccess),
	})
}

func (h *AdminAnalysisHandler) UpdateAnalysis(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var updateInput models.AdminAnalysisUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &updateInput); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if updateInput.Status != nil && !updateInput.Status.IsValid() {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.AnalysisInvalidStatus),
		})
		return
	}

	analysisUpdated, err := h.Service.Update(c.Request.Context(),
		id, updateInput)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: analysisUpdated})
}

func (h *AdminAnalysisHandler) DeleteAnalysis(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.AnalysisDeleted),
	})
}

func (h *AdminAnalysisHandler) DownloadZip(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	analysis, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	if analysis.ResultsZipPath == nil {
		c.JSON(http.StatusNotFound, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.AnalysisZipNotFound),
		})
		return
	}

	c.File(*analysis.ResultsZipPath)
}

func (h *AdminAnalysisHandler) DownloadBatchTSV(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var downloadInput models.AnalysisTSVDownloadInput
	if errMsg, valid := validations.Validate(c, localizer,
		&downloadInput); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	analyses, err := h.Service.FindManyByIDs(c.Request.Context(),
		downloadInput.IDs)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	tsvBytes, err := utils.GenerateDynamicTSV(analyses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.GenericInternalServerError),
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=cabgen_results.tsv")
	c.Data(http.StatusOK, "text/tab-separated-values", tsvBytes)
}
