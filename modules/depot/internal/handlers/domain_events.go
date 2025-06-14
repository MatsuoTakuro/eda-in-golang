package handlers

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/modules/depot/internal/domain"
)

func SubscribeDomainEvents(container di.Container) {

	handler := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {

		domainHandlers := di.Get(ctx, di.DomainEventHandler).(ddd.EventHandler[ddd.AggregateEvent])

		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(di.DomainDispatcher).(ddd.EventDispatcher[ddd.AggregateEvent])

	subscribeDomainEvents(subscriber, handler)
}

func subscribeDomainEvents(
	subscriber ddd.EventSubscriber[ddd.AggregateEvent],
	handler ddd.EventHandler[ddd.AggregateEvent],
) {
	subscriber.Subscribe(handler,
		domain.ShoppingListCompletedEvent,
	)
}
