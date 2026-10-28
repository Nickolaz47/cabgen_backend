package container

import (
	adminHandler "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/ticket"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/contact"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildTicketService(db *gorm.DB, taskQueue services.TaskQueue,
	logger *zap.Logger) services.TicketService {
	ticketRepo := repositories.NewTicketRepo(db)
	ticketService := services.NewTicketService(
		ticketRepo, taskQueue, logger,
	)

	return ticketService
}

func BuildTicketHandler(svc services.TicketService,
) *contact.TicketHandler {
	return contact.NewTicketHandler(svc)
}

func BuildAdminTicketHandler(svc services.TicketService,
) *adminHandler.AdminTicketHandler {
	return adminHandler.NewAdminTicketHandler(svc)
}
