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

type AdminAnalysisHandler struct {
	Service services.AnalysisService
}

func NewAdminAnalysisHandler(svc services.AnalysisService,
) *AdminAnalysisHandler {
	return &AdminAnalysisHandler{
		Service: svc,
	}
}

func (h *AdminAnalysisHandler) GetAnalyses(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminAnalysesGet, nil)
	language := translation.GetLanguageFromContext(c)

	var filter models.AnalysisFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminAnalysesGetFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.InvalidQueryParamError),
		})
		return
	}

	analyses, err := h.Service.FindAll(c.Request.Context(), uuid.Nil,
		filter, language)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminAnalysesGetFailed, nil)
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
	validations.SetAuditEvent(c, models.AuditEventAdminAnalysesGetByID, nil)
	language := translation.GetLanguageFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminAnalysesGetByIDFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	analysis, err := h.Service.FindByID(c.Request.Context(), id, uuid.Nil, language)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminAnalysesGetByIDFailed, nil)
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: analysis})
}

func (h *AdminAnalysisHandler) CreateAnalysis(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c, models.AuditEventAdminAnalysesCreate, nil)
	language := translation.GetLanguageFromContext(c)

	var newAnalysis models.AdminAnalysisCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &newAnalysis); !valid {
		validations.SetAuditEvent(c, models.AuditEventAdminAnalysesCreateFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if !newAnalysis.Type.IsValid() {
		validations.SetAuditEvent(c, models.AuditEventAdminAnalysesCreateFailed, nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.AnalysisInvalidType),
		})
		return
	}

	payload := models.AnalysisCreateDTO(newAnalysis)
	analysis, err := h.Service.Create(c.Request.Context(), payload, language)
	if err != nil {
		validations.SetAuditEvent(c, models.AuditEventAdminAnalysesCreateFailed, nil)
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
	validations.SetAuditEvent(c, models.AuditEventAdminAnalysesUpdate, nil)
	language := translation.GetLanguageFromContext(c)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesUpdateFailed,
			map[string]string{"analysis_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var updateInput models.AdminAnalysisUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &updateInput); !valid {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesUpdateFailed,
			map[string]string{"analysis_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if updateInput.Status != nil && !updateInput.Status.IsValid() {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesUpdateFailed,
			map[string]string{"analysis_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.AnalysisInvalidStatus),
		})
		return
	}

	analysisUpdated, err := h.Service.Update(c.Request.Context(),
		id, updateInput, language)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesUpdateFailed,
			map[string]string{"analysis_id": rawID})
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
	validations.SetAuditEvent(c, models.AuditEventAdminAnalysesDelete, nil)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesDeleteFailed,
			map[string]string{"analysis_id": rawID})
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id, uuid.Nil); err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesDeleteFailed,
			map[string]string{"analysis_id": rawID})
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
	validations.SetAuditEvent(c, models.AuditEventAdminAnalysesDownloadZip, nil)
	rawID := c.Param("analysisId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesDownloadZipFailed,
			nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	zipPath, err := h.Service.DownloadZip(c.Request.Context(), id, uuid.Nil)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesDownloadZipFailed,
			nil)
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

func (h *AdminAnalysisHandler) DownloadBatchTSV(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	validations.SetAuditEvent(c,
		models.AuditEventAdminAnalysesDownloadBatchTSV,
		nil)
	language := translation.GetLanguageFromContext(c)

	var downloadInput models.AnalysisTSVDownloadInput
	if errMsg, valid := validations.Validate(c, localizer,
		&downloadInput); !valid {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesDownloadBatchTSVFailed,
			nil)
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	analyses, err := h.Service.DownloadBatchTSV(c.Request.Context(),
		downloadInput.IDs, uuid.Nil, language)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesDownloadBatchTSVFailed,
			nil)
		code, errMsg := handlererrors.HandleAnalysisError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	tsvBytes, err := utils.GenerateMetricsTSV(analyses)
	if err != nil {
		validations.SetAuditEvent(c,
			models.AuditEventAdminAnalysesDownloadBatchTSVFailed,
			nil)
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.GenericInternalServerError),
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=cabgen_results.tsv")
	c.Data(http.StatusOK, "text/tab-separated-values", tsvBytes)
}
