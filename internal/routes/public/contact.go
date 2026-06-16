package public

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/contact"
	"github.com/gin-gonic/gin"
)

func SetupContactRoutes(r *gin.RouterGroup, handler *contact.TicketHandler) {
	contactRouter := r.Group("/contact")

	contactRouter.POST("", handler.CreateTicket)
}
