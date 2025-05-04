package eventhandlers

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type Invoice interface {
	OnOrderReadiedEventHandler
}

type invoice struct {
	client infra.InvoiceClient
}

var _ Invoice = (*invoice)(nil)

func NewInvoice(client infra.InvoiceClient) Invoice {
	return &invoice{
		client: client,
	}
}

func (h invoice) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	orderReadied := event.(*domain.OrderReadied)
	return h.client.Save(ctx, orderReadied.Order.ID, orderReadied.Order.PaymentID, orderReadied.Order.GetTotal())
}
