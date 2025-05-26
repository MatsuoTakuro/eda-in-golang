package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/payments/internal/domain"
)

func SubscribeDomainEventsForIntegration(eventHandlers ddd.EventHandler[ddd.Event], domainSubscriber ddd.EventSubscriber[ddd.Event]) {
	domainSubscriber.Subscribe(eventHandlers,
		domain.InvoicePaidEvent,
	)
}
