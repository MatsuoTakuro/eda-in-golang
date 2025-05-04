package ddd

import (
	"context"
	"sync"
)

type EventDispatcher interface {
	EventSubscriber
	EventPublisher
}

type EventSubscriber interface {
	Subscribe(event Event, handler EventHandler)
}

type EventPublisher interface {
	Publish(ctx context.Context, events ...Event) error
}

type eventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.Mutex
}

var _ EventDispatcher = (*eventDispatcher)(nil)

func NewEventDispatcher() EventDispatcher {
	return &eventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (h *eventDispatcher) Subscribe(event Event, handler EventHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.handlers[event.EventName()] = append(h.handlers[event.EventName()], handler)
}

func (h *eventDispatcher) Publish(ctx context.Context, events ...Event) error {
	for _, event := range events {
		for _, handler := range h.handlers[event.EventName()] {
			err := handler(ctx, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
