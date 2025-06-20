package application

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/depot/internal/domain"
)

type domainHandler[T ddd.AggregateEvent] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*domainHandler[ddd.AggregateEvent])(nil)

// NewDomainEventHandler creates a handler that handle domain events occuring within the module.
func NewDomainEventHandler(publisher am.EventPublisher) domainHandler[ddd.AggregateEvent] {
	return domainHandler[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ShoppingListCompletedEvent:
		return h.onShoppingListCompleted(ctx, event)
	}
	return nil
}

func (h domainHandler[T]) onShoppingListCompleted(ctx context.Context, event ddd.AggregateEvent) error {
	completed := event.Payload().(*domain.ShoppingListCompleted)

	return h.publisher.Publish(ctx,
		depotpb.ShoppingListAggregateChannel,
		ddd.NewEvent(depotpb.ShoppingListCompletedEvent,
			&depotpb.ShoppingListCompleted{
				Id:      event.AggregateID(),
				OrderId: completed.ShoppingList.OrderID,
			},
		),
	)
}
