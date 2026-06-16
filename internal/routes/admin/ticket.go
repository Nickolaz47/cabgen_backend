package admin

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/ticket"
	"github.com/gin-gonic/gin"
)

func SetupAdminTicketRoutes(r *gin.RouterGroup,
	handler *ticket.AdminTicketHandler) {
	ticketRouter := r.Group("/tickets")

	ticketRouter.GET("", handler.GetTickets)
	ticketRouter.GET("/:ticketId", handler.GetTicketByID)
	ticketRouter.PUT("/:ticketId/assign", handler.Assign)
	ticketRouter.PUT("/:ticketId/resolve", handler.Resolve)
	ticketRouter.DELETE("/:ticketId", handler.DeleteTicket)
}
