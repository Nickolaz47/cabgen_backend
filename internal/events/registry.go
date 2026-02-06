package events

import (
	"context"
	"fmt"
)

type HandlerFunc func(ctx context.Context, payload []byte) error

type Registry interface {
	Register(eventName string, handler HandlerFunc) error
	Get(eventName string) (HandlerFunc, bool)
}

type registry struct {
	handlers map[string]HandlerFunc
}

func NewRegistry() Registry {
	return &registry{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *registry) Register(eventName string, handler HandlerFunc) error {
	if _, exists := r.handlers[eventName]; exists {
		return fmt.Errorf("handler already registered for event %s", eventName)
	}
	r.handlers[eventName] = handler

	return nil
}

func (r *registry) Get(eventName string) (HandlerFunc, bool) {
	h, ok := r.handlers[eventName]
	return h, ok
}
