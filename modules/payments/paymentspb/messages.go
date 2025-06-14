package paymentspb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	InvoiceAggregateChannel = "mallbots.payments.events.Invoice"

	InvoicePaidEvent = "paymentsapi.InvoicePaid"

	CommandChannel = "mallbots.payments.commands"

	ConfirmPaymentCommand = "paymentsapi.ConfirmPayment"
)

func Registrations(reg registry.Registry) (err error) {
	regtr := registrar.NewProtoRegistrar(reg)

	// Invoice events
	if err = regtr.Register(&InvoicePaid{}); err != nil {
		return err
	}

	// commands
	if err = regtr.Register(&ConfirmPayment{}); err != nil {
		return
	}

	return
}

func (*InvoicePaid) Key() string { return InvoicePaidEvent }

func (*ConfirmPayment) Key() string { return ConfirmPaymentCommand }
