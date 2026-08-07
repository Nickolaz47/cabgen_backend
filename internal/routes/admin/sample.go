package admin

import (
	sample "github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/gin-gonic/gin"
)

func SetupAdminSampleRoutes(r *gin.RouterGroup,
	handler *sample.SampleHandler) {
	sampleRouter := r.Group("/samples")

	sampleRouter.GET("", handler.GetSamples)
	sampleRouter.GET("/:sampleId", handler.GetSampleByID)
	sampleRouter.POST("", handler.CreateSample)
	sampleRouter.PUT("/:sampleId/upload", handler.UploadFiles)
	sampleRouter.PUT("/:sampleId", handler.UpdateSample)
	sampleRouter.DELETE("/:sampleId", handler.DeleteSample)
}
