package origin

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func GetActiveOrigins(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	origins, err := repository.OriginRepo.GetActiveOrigins()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	formOrigins := make([]models.OriginFormResponse, len(origins))
	for i, orig := range origins {
		formOrigins[i] = orig.ToFormResponse(c)
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: formOrigins,
	})
}
