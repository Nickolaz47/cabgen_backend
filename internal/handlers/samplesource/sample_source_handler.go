package samplesource

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func GetActiveSampleSources(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	sampleSources, err := repository.SampleSourceRepo.GetActiveSampleSources()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}
	
	publicSampleSources := make([]models.SampleSourceFormResponse, len(sampleSources))
	for i, s := range sampleSources {
		publicSampleSources[i] = s.ToFormResponse(c)
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: publicSampleSources,
	})
}
