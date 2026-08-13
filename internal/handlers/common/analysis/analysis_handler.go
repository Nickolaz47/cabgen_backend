package analysis

import (
	"net/http"
	"path/filepath"

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

func (h *AnalysisHandler) GetAnalysisFastQCByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		handlererrors.RespondHTMLError(c, http.StatusBadRequest,
			responses.GetResponse(localizer, responses.InvalidURLID))
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		handlererrors.RespondHTMLError(c, http.StatusUnauthorized,
			responses.GetResponse(localizer, responses.UnauthorizedError))
		return
	}

	analysis, err := h.Service.FindByID(c.Request.Context(), id, userToken.ID)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		handlererrors.RespondHTMLError(c, code,
			responses.GetResponse(localizer, errMsg))
		return
	}

	var htmlPath *string
	requestedFastqc := c.Param("fastqcReport")

	switch requestedFastqc {
	case "fastqc1":
		htmlPath = analysis.FastQC1
	case "fastqc2":
		htmlPath = analysis.FastQC2
	default:
		handlererrors.RespondHTMLError(c, http.StatusNotFound,
			responses.GetResponse(localizer,
				responses.AnalysisInvalidFastQCReport))
		return
	}

	if htmlPath == nil {
		handlererrors.RespondHTMLError(c, http.StatusNotFound,
			responses.GetResponse(localizer,
				responses.AnalysisFastQCReportNotAvailable))
		return
	}

	c.File(*htmlPath)
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
		Data: analysis,
		Message: responses.GetResponse(localizer,
			responses.AnalysisCreationSuccess),
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

	zipPath, err := h.Service.DownloadZip(c.Request.Context(), id,
		userToken.ID)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.Header("Content-Disposition",
		"attachment; filename="+filepath.Base(zipPath))
	c.File(zipPath)
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

	analyses, err := h.Service.DownloadBatchTSV(c.Request.Context(),
		downloadInput.IDs, userToken.ID)
	if err != nil {
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	tsvBytes, err := utils.GenerateMetricsTSV(analyses)
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
