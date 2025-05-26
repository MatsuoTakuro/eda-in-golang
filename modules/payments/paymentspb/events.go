package paymentspb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	InvoiceAggregateChannel = "mallbots.payments.events.Invoice"

	InvoicePaidEvent = "paymentsapi.InvoicePaid"
)

func RegisterIntegrationEvents(reg registry.Registry) error {
	regtr := registrar.NewProtoRegistrar(reg)

	// Invoice events
	if err := regtr.Register(&InvoicePaid{}); err != nil {
		return err
	}

	return nil
}

func (*InvoicePaid) Key() string { return InvoicePaidEvent }
