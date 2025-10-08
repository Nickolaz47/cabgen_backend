package origin

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func GetAllActiveOrigins(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	origins, err := repository.OriginRepo.GetActiveOrigins()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	publicOrigins := make([]models.OriginPublicResponse, len(origins))
	for i, orig := range origins {
		publicOrigins[i] = orig.ToPublicResponse(c)
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: publicOrigins,
	})
}
