package sequencer

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

func GetAllSequencers(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	sequencers, err := repository.SequencerRepo.GetSequencers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sequencers})
}

func GetSequencerByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sequencerId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	sequencer, err := repository.SequencerRepo.GetSequencerByID(c.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.SequencerNotFoundError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sequencer})
}

func GetSequencersByBrandOrModel(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	input := c.Query("brandOrModel")

	if input == "" {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.SequencerEmptyQueryError),
		})
		return
	}

	sequencers, err := repository.SequencerRepo.GetSequencersByBrandOrModel(c.Request.Context(),input)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: sequencers})
}

func CreateSequencer(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newSequencer models.SequencerCreateInput
	if errMsg, valid := validations.Validate(c, localizer, &newSequencer); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	sequencerToCreate := models.Sequencer{
		Brand:    newSequencer.Brand,
		Model:    newSequencer.Model,
		IsActive: newSequencer.IsActive,
	}

	if err := repository.SequencerRepo.CreateSequencer(c.Request.Context(), &sequencerToCreate); err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.SequencerCreationSuccess),
		Data:    sequencerToCreate,
	})
}

func UpdateSequencer(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sequencerId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var sequencerUpdateInput models.SequencerUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &sequencerUpdateInput); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	sequencerToUpdate, err := repository.SequencerRepo.GetSequencerByID(c.Request.Context(),id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.SequencerNotFoundError),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	validations.ApplySequencerUpdate(sequencerToUpdate, &sequencerUpdateInput)
	if err := repository.SequencerRepo.UpdateSequencer(c.Request.Context(), sequencerToUpdate); err != nil {
		c.JSON(http.StatusOK, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: sequencerToUpdate,
	})
}

func DeleteSequencer(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("sequencerId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	sequencerToDelete, err := repository.SequencerRepo.GetSequencerByID(c.Request.Context(),id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.SequencerNotFoundError),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	if err := repository.SequencerRepo.DeleteSequencer(c.Request.Context(), sequencerToDelete); err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.SequencerDeleted),
	})
}
