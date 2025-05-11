package eventhandlers

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type invoice[T ddd.AggregateEvent] struct {
	client infra.InvoiceClient
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*invoice[ddd.AggregateEvent])(nil)

func NewInvoice(client infra.InvoiceClient) *invoice[ddd.AggregateEvent] {
	return &invoice[ddd.AggregateEvent]{
		client: client,
	}
}

func (h *invoice[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	}
	return nil
}

func (h invoice[T]) onOrderReadied(ctx context.Context, event ddd.AggregateEvent) error {
	orderReadied := event.Payload().(*domain.OrderReadied)
	return h.client.Save(ctx, event.AggregateID(), orderReadied.PaymentID, orderReadied.Total)
}
