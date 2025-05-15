package eventhandlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
)

func SubscribeForInvoice(invoiceHandler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.OrderReadiedEvent, invoiceHandler)
}

func SubscribeForNotification(notificationHandler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.OrderCreatedEvent, notificationHandler)
	domainSubscriber.Subscribe(domain.OrderReadiedEvent, notificationHandler)
	domainSubscriber.Subscribe(domain.OrderCanceledEvent, notificationHandler)
}
