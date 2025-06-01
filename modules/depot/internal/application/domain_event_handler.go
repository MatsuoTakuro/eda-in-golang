package application

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/internal/domain"
)

type domainHandlers[T ddd.AggregateEvent] struct {
	orders domain.OrderRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*domainHandlers[ddd.AggregateEvent])(nil)

// NewDomainEventHandlers creates a handler that handle domain events occuring within the module.
func NewDomainEventHandlers(
	orders domain.OrderRepository,
) domainHandlers[ddd.AggregateEvent] {
	return domainHandlers[ddd.AggregateEvent]{
		orders: orders,
	}
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ShoppingListCompletedEvent:
		return h.onShoppingListCompleted(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onShoppingListCompleted(ctx context.Context, event ddd.AggregateEvent) error {
	completed := event.Payload().(*domain.ShoppingListCompleted)
	return h.orders.Ready(ctx, completed.ShoppingList.OrderID)
}
