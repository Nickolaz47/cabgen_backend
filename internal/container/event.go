package container

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/events"
	"github.com/CABGenOrg/cabgen_backend/internal/events/handlers"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildEventRepository(db *gorm.DB) repositories.EventRepository {
	return repositories.NewEventRepo(db)
}

func BuildRegistry(emailSvc services.EmailService) events.Registry {
	handlersToRegister := map[string]events.HandlerFunc{
		"user.registered": handlers.UserRegisteredHandler(emailSvc),
	}

	registry := events.NewRegistry()
	for name, handler := range handlersToRegister {
		registry.Register(name, handler)
	}

	return registry
}

func BuildEventDispatcher(
	eventRepo repositories.EventRepository,
	registry events.Registry) events.Dispatcher {
	interval := 5 * time.Minute
	nWorkers := 5

	return events.NewDispatcher(eventRepo, registry, interval, nWorkers)
}
