package eventhandlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
)

func SubscribeDomainEventsForInvoice(invoiceHandler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(invoiceHandler, domain.OrderReadiedEvent)
}

func SubscribeDomainEventsForNotification(notificationHandler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(notificationHandler,
		domain.OrderCreatedEvent,
		domain.OrderReadiedEvent,
		domain.OrderCanceledEvent,
	)
}
