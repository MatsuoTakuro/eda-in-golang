package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/internal/domain"
)

func SubscribeDomainEventsForOrder(orderHandler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(orderHandler, domain.ShoppingListCompletedEvent)
}
