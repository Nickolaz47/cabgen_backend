package sequencer

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func GetActiveSequencers(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	activeSequencers, err := repository.SequencerRepo.GetActiveSequencers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.GenericInternalServerError),
		})
		return
	}

	formSequencers := make([]models.SequencerFormResponse, len(activeSequencers))
	for i, seq := range activeSequencers {
		formSequencers[i] = seq.ToFormResponse()
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: formSequencers})
}
