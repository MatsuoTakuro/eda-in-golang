package eventhandlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
)

func SubscribeForInvoice(invoiceHandler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.OrderReadiedEvent, invoiceHandler)
}
