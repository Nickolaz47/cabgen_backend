package origin

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

func GetAllOrigins(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	origins, err := repository.OriginRepo.GetOrigins()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: origins})
}

func GetOriginByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("originId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	origin, err := repository.OriginRepo.GetOriginByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.OriginNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: origin})
}

func GetOriginByName(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.OriginEmptyNameError),
		})
		return
	}

	origin, err := repository.OriginRepo.GetOriginByName(name, language)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.OriginNotFoundError)},
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
		Data: origin,
	})
}

func CreateOrigin(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newOrigin models.OriginCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &newOrigin); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	errMsg, ok := validations.ValidateTranslationMap(c, "origin", newOrigin.Names)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	originToCreate := models.Origin{
		Names:    newOrigin.Names,
		IsActive: newOrigin.IsActive,
	}

	if err := repository.OriginRepo.CreateOrigin(&originToCreate); err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data:    originToCreate.ToResponse(c),
		Message: responses.GetResponse(localizer, responses.OriginCreationSuccess),
	})
}

func UpdateOrigin(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("originId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var originUpdateInput models.OriginUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &originUpdateInput)
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	if originUpdateInput.Names != nil {
		errMsg, ok = validations.ValidateTranslationMap(c, "origin", originUpdateInput.Names)
	}

	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	originToUpdate, err := repository.OriginRepo.GetOriginByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.OriginNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	validations.ApplyOriginUpdate(originToUpdate, &originUpdateInput)
	if err := repository.OriginRepo.UpdateOrigin(originToUpdate); err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(
				localizer, responses.GenericInternalServerError,
			)})
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Data: originToUpdate.ToResponse(c),
		},
	)
}

func DeleteOrigin(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("originId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	origin, err := repository.OriginRepo.GetOriginByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.OriginNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err := repository.OriginRepo.DeleteOrigin(origin); err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Message: responses.GetResponse(localizer, responses.OriginDeleted),
		})
}
