package eventhandlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/application/eventhandlers"
	"eda-in-golang/modules/ordering/internal/domain"
)

func SubscribeForInvoice(invoice eventhandlers.Invoice, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderReadied{}, invoice.OnOrderReadied)
}
