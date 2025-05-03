package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/baskets/internal/application"
	"eda-in-golang/modules/baskets/internal/domain"
)

func RegisterOrderHandlers(orderHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.BasketCheckedOut{}, orderHandlers.OnBasketCheckedOut)
}
