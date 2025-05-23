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
	Subscribe(handler EventHandler[T], events ...string)
}

// EventPublisher defines how to publish one or more events to registered handlers.
type EventPublisher[T Event] interface {
	Publish(ctx context.Context, events ...T) error
}

type eventDispatcher[T Event] struct {
	handlers []eventHandleFilters[T]
	mu       sync.Mutex
}

type eventHandleFilters[T Event] struct {
	h       EventHandler[T]
	filters map[string]struct{}
}

var _ EventDispatcher[Event] = (*eventDispatcher[Event])(nil)

// NewEventDispatcher creates a new domain event dispatcher instance.
func NewEventDispatcher[T Event]() *eventDispatcher[T] {
	return &eventDispatcher[T]{
		handlers: make([]eventHandleFilters[T], 0),
	}
}

func (ed *eventDispatcher[T]) Subscribe(handler EventHandler[T], events ...string) {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	var filters map[string]struct{}
	if len(events) > 0 {
		filters = make(map[string]struct{})
		for _, event := range events {
			filters[event] = struct{}{}
		}
	}

	ed.handlers = append(ed.handlers, eventHandleFilters[T]{
		h:       handler,
		filters: filters,
	})
}

func (ed *eventDispatcher[T]) Publish(ctx context.Context, events ...T) error {
	for _, event := range events {
		for _, handler := range ed.handlers {
			if handler.filters != nil {
				if _, exists := handler.filters[event.EventName()]; !exists {
					continue
				}
			}
			err := handler.h.HandleEvent(ctx, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
