package samplesource

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSampleSources(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	sampleSources, err := repository.SampleSourceRepo.GetSampleSources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sampleSources})
}

func GetSampleSourceByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleSourceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	sampleSource, err := repository.SampleSourceRepo.GetSampleSourceByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.SampleSourceNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: sampleSource,
	})
}

func GetSampleSourceByNameOrGroup(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	input := c.Query("nameOrGroup")

	if input == "" {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.SampleSourceEmptyQueryError),
		})
		return
	}

	sampleSources, err := repository.SampleSourceRepo.GetSampleSourcesByNameOrGroup(input, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: sampleSources,
	})
}

func CreateSampleSource(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newSampleSource models.SampleSourceCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &newSampleSource); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	errMsg, ok := validations.ValidateTranslationMap(c, "sampleSource", newSampleSource.Names)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	errMsg, ok = validations.ValidateTranslationMap(c, "sampleSource", newSampleSource.Groups)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	sampleSourceToCreate := models.SampleSource{
		Names:    newSampleSource.Names,
		Groups:   newSampleSource.Groups,
		IsActive: newSampleSource.IsActive,
	}

	if err := repository.SampleSourceRepo.CreateSampleSource(&sampleSourceToCreate); err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.SampleSourceCreationSuccess),
		Data:    sampleSourceToCreate.ToResponse(c),
	})
}

func UpdateSampleSource(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleSourceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var sampleSourceUpdateInput models.SampleSourceUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &sampleSourceUpdateInput)
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	if sampleSourceUpdateInput.Names != nil {
		errMsg, ok = validations.ValidateTranslationMap(c, "sampleSource", sampleSourceUpdateInput.Names)
	}

	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	if sampleSourceUpdateInput.Groups != nil {
		errMsg, ok = validations.ValidateTranslationMap(c, "sampleSource", sampleSourceUpdateInput.Groups)
	}

	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	sampleSourceToUpdate, err := repository.SampleSourceRepo.GetSampleSourceByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.SampleSourceNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	validations.ApplySampleSourceUpdate(sampleSourceToUpdate, &sampleSourceUpdateInput)

	if err := repository.SampleSourceRepo.UpdateSampleSource(sampleSourceToUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: sampleSourceToUpdate.ToResponse(c),
	})
}

func DeleteSampleSource(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sampleSourceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	sampleSourceToDelete, err := repository.SampleSourceRepo.GetSampleSourceByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.SampleSourceNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err := repository.SampleSourceRepo.DeleteSampleSource(sampleSourceToDelete); err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.SampleSourceDeleted),
	})
}
