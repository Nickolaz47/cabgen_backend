package admin

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
	rawID := c.Param("originID")
	id := uuid.MustParse(rawID)

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
	name := c.Query("originName")

	if name == "" {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.OriginNameEmptyError),
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

	errMsg, ok := validations.ValidateOriginNames(c, &newOrigin)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
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

func UpdateOrigin(c *gin.Context) {}

func DeleteOrigin(c *gin.Context) {}
