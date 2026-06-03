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

type AnalysisHandler struct {
	Service services.AnalysisService
}

func NewAnalysisHandler(svc services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{
		Service: svc,
	}
}

func (h *AnalysisHandler) GetAnalyses(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UnauthorizedError),
		})
		return
	}

	analyses, err := h.Service.FindAll(c.Request.Context(), userToken.ID)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: analyses})
}

func (h *AnalysisHandler) GetAnalysisByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UnauthorizedError),
		})
		return
	}

	analysis, err := h.Service.FindByID(c.Request.Context(), id, userToken.ID)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: analysis})
}

func (h *AnalysisHandler) CreateAnalysis(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newAnalysis models.AnalysisCreateInput
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

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.UnauthorizedError),
		})
		return
	}

	payload := models.AnalysisCreateInputToDTO(newAnalysis, userToken.ID)
	analysis, err := h.Service.Create(c.Request.Context(), payload)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data:    analysis,
		Message: responses.GetResponse(localizer, responses.AnalysisCreationSuccess),
	})
}

func (h *AnalysisHandler) DeleteAnalysis(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UnauthorizedError),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id, userToken.ID); err != nil {
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

func (h *AnalysisHandler) DownloadZip(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.UnauthorizedError),
		})
		return
	}

	analysis, err := h.Service.FindByID(c.Request.Context(), id, userToken.ID)
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

func (h *AnalysisHandler) DownloadBatchTSV(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var downloadInput models.AnalysisTSVDownloadInput
	if errMsg, valid := validations.Validate(c, localizer,
		&downloadInput); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.UnauthorizedError),
		})
		return
	}

	analyses, err := h.Service.FindManyByIDs(c.Request.Context(),
		downloadInput.IDs, userToken.ID)
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
