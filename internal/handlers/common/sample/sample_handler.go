package sample

import (
	"io"
	"net/http"
	"os"
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

type Scope int

const (
	ScopeSelf Scope = iota // filters by userID
	ScopeAll               // no filter, returns all data
)

type SampleHandler struct {
	Service services.SampleService
	Scope   Scope
}

func NewSampleHandler(svc services.SampleService) *SampleHandler {
	return &SampleHandler{
		Service: svc,
		Scope:   ScopeSelf,
	}
}

func NewAdminSampleHandler(svc services.SampleService) *SampleHandler {
	return &SampleHandler{
		Service: svc,
		Scope:   ScopeAll,
	}
}

func (h *SampleHandler) getUserID(userToken *models.UserToken) uuid.UUID {
	if h.Scope == ScopeAll {
		return uuid.Nil
	}
	return userToken.ID
}

func (h *SampleHandler) GetSamples(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	input := utils.SanitizeQuery(c.Query("input"))

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	samples, err := h.Service.FindAll(c.Request.Context(), input,
		h.getUserID(userToken), language)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: samples})
}

func (h *SampleHandler) GetSampleByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	rawID := c.Param("sampleId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	sample, err := h.Service.FindByID(c.Request.Context(), id,
		h.getUserID(userToken), language)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sample})
}

func (h *SampleHandler) CreateSample(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var payload models.SampleCreateDTO
	if h.Scope == ScopeAll {
		var newSample models.AdminSampleCreateInput
		if errMsg, valid := validations.Validate(c, localizer, &newSample); !valid {
			c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
			return
		}

		if newSample.Gender != nil && !newSample.Gender.IsValid() {
			c.JSON(http.StatusBadRequest,
				responses.APIResponse{
					Error: responses.GetResponse(localizer,
						responses.SampleInvalidGender),
				})
			return
		}

		payload = models.SampleCreateDTO(newSample)
	} else {
		var newSample models.SampleCreateInput
		if errMsg, valid := validations.Validate(c, localizer, &newSample); !valid {
			c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
			return
		}

		if newSample.Gender != nil && !newSample.Gender.IsValid() {
			c.JSON(http.StatusBadRequest,
				responses.APIResponse{
					Error: responses.GetResponse(localizer,
						responses.SampleInvalidGender),
				})
			return
		}

		payload = models.SampleCreateInputToDTO(newSample, userToken.ID)
	}

	sample, err := h.Service.Create(c.Request.Context(), payload, language)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleError(err)
		c.JSON(code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data: sample,
		Message: responses.GetResponse(localizer,
			responses.SampleCreationSuccess),
	})
}

func (h *SampleHandler) UploadFiles(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer,
				responses.SampleContentTypeError),
		})
		return
	}

	uploadDir, err := h.Service.PrepareSampleFolder(userToken.ID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.GenericInternalServerError)})
		return
	}

	var attachmentInput models.SampleAttachmentInput

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			// Finish the request
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.GenericInternalServerError),
			})
			return
		}

		formName := part.FormName()
		fileName := filepath.Base(part.FileName())

		// Skip form parts that are not files
		if fileName == "" || fileName == "." {
			continue
		}

		// Save to the disk
		dstPath := filepath.Join(uploadDir, fileName)

		out, err := os.Create(dstPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.APIResponse{
				Error: responses.GetResponse(
					localizer, responses.GenericInternalServerError,
				),
			})
			return
		}

		// Do the streaming from net to disk
		_, err = io.Copy(out, part)
		out.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.GenericInternalServerError),
			})
			return
		}

		switch formName {
		case "fastq1":
			attachmentInput.Fastq1 = &fileName
		case "fastq2":
			attachmentInput.Fastq2 = &fileName
		case "fasta":
			attachmentInput.Fasta = &fileName
		}
	}

	if err := h.Service.AttachFiles(c.Request.Context(),
		id, h.getUserID(userToken), attachmentInput); err != nil {
		code, errMsg := handlererrors.HandleSampleError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer,
			responses.SampleUploadSuccess),
	})
}

func (h *SampleHandler) UpdateSample(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	rawID := c.Param("sampleId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var payload models.SampleUpdateDTO
	if h.Scope == ScopeAll {
		var sampleUpdateInput models.AdminSampleUpdateInput
		errMsg, ok := validations.Validate(c, localizer, &sampleUpdateInput)
		if !ok {
			c.JSON(http.StatusBadRequest,
				responses.APIResponse{
					Error: errMsg,
				})
			return
		}

		if sampleUpdateInput.Gender != nil &&
			!sampleUpdateInput.Gender.IsValid() {
			c.JSON(http.StatusBadRequest, responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.SampleInvalidGender),
			})
			return
		}

		payload = models.SampleUpdateDTO(sampleUpdateInput)
	} else {
		var sampleUpdateInput models.SampleUpdateInput
		errMsg, ok := validations.Validate(c, localizer, &sampleUpdateInput)
		if !ok {
			c.JSON(http.StatusBadRequest,
				responses.APIResponse{
					Error: errMsg,
				})
			return
		}

		if sampleUpdateInput.Gender != nil &&
			!sampleUpdateInput.Gender.IsValid() {
			c.JSON(http.StatusBadRequest, responses.APIResponse{
				Error: responses.GetResponse(localizer,
					responses.SampleInvalidGender),
			})
			return
		}

		payload = models.SampleUpdateInputToDTO(sampleUpdateInput,
			userToken.ID)
	}

	sampleUpdated, err := h.Service.Update(c.Request.Context(), id,
		h.getUserID(userToken), payload, language)
	if err != nil {
		code, errMsg := handlererrors.HandleSampleError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sampleUpdated})
}

func (h *SampleHandler) DeleteSample(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id,
		h.getUserID(userToken)); err != nil {
		code, errMsg := handlererrors.HandleSampleError(err)
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
			localizer, responses.SampleDeleted,
		)},
	)
}
