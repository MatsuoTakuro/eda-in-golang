package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/internal/application"
	"eda-in-golang/modules/depot/internal/domain"
)

func RegisterOrderHandlers(orderHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.ShoppingListCompleted{}, orderHandlers.OnShoppingListCompleted)
}
