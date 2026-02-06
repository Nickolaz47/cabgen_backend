package events

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
)

type EventEmitter interface {
	Emit(ctx context.Context, name string, payload any) error
}

type eventEmitter struct {
	repo repositories.EventRepository
}

func NewEventEmitter(repo repositories.EventRepository) EventEmitter {
	return &eventEmitter{
		repo: repo,
	}
}

func (e *eventEmitter) Emit(
	ctx context.Context, name string, payload any) error {
	ev, err := models.NewEvent(name, payload)
	if err != nil {
		return err
	}

	return e.repo.CreateEvent(ctx, ev)
}
