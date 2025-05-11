package ddd

import (
	"context"
	"sync"
)

// EventDispatcher combines subscription and publication capabilities for domain events.
type EventDispatcher[T Event] interface {
	EventSubscriber[T]
	EventPublisher[T]
}

// EventSubscriber allows registering event handlers by event name.
type EventSubscriber[T Event] interface {
	Subscribe(name string, handler EventHandler[T])
}

// EventPublisher defines how to publish one or more events to registered handlers.
type EventPublisher[T Event] interface {
	Publish(ctx context.Context, events ...T) error
}

type eventDispatcher[T Event] struct {
	handlers map[string][]EventHandler[T]
	mu       sync.Mutex
}

var _ EventDispatcher[Event] = (*eventDispatcher[Event])(nil)

// NewEventDispatcher creates a new domain event dispatcher instance.
func NewEventDispatcher[T Event]() *eventDispatcher[T] {
	return &eventDispatcher[T]{
		handlers: make(map[string][]EventHandler[T]),
	}
}

func (d *eventDispatcher[T]) Subscribe(name string, handler EventHandler[T]) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers[name] = append(d.handlers[name], handler)
}

func (d *eventDispatcher[T]) Publish(ctx context.Context, events ...T) error {
	for _, event := range events {
		// For now, synchronously invoke all handlers registered for this event
		for _, h := range d.handlers[event.EventName()] {
			err := h.HandleEvent(ctx, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
