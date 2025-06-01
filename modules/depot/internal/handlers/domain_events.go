package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/internal/domain"
)

func SubscribeDomainEvents(
	subscriber ddd.EventSubscriber[ddd.AggregateEvent],
	handler ddd.EventHandler[ddd.AggregateEvent],
) {
	subscriber.Subscribe(handler, domain.ShoppingListCompletedEvent)
}
