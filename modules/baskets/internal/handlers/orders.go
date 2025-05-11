package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/baskets/internal/domain"
)

func RegisterOrderHandlers(orderHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.BasketCheckedOutEvent, orderHandlers)
}
