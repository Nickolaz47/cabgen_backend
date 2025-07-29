package routes

import (
	"github.com/CABGenOrg/cabgen_backend/internal/routes/country"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/public"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	public.PublicRoutes(router)
	country.CountryRoutes(router)
}
