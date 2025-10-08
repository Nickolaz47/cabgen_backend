package routes

import (
	"github.com/CABGenOrg/cabgen_backend/internal/routes/admin"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/country"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/origin"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/public"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/user"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	public.PublicRoutes(router)
	country.CountryRoutes(router)
	user.UserRoutes(router)
	origin.OriginRoutes(router)
	admin.AdminRoutes(router)
}
