package public

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	c.JSON(http.StatusOK,
		responses.APIResponse{
			Message: responses.GetResponse(localizer, responses.HealthMessage),
		},
	)
}
