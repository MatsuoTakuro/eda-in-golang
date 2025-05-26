package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/customers/internal/domain"
)

func SubscribeDomainEventsForIntegration(eventHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(eventHandlers,
		domain.CustomerRegisteredEvent,
		domain.CustomerSmsChangedEvent,
		domain.CustomerEnabledEvent,
		domain.CustomerDisabledEvent,
	)
}
