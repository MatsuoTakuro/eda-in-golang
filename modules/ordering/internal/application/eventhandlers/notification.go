package eventhandlers

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type notification[T ddd.AggregateEvent] struct {
	client infra.NotificationClient
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*notification[ddd.AggregateEvent])(nil)

func NewNotification(client infra.NotificationClient) *notification[ddd.AggregateEvent] {
	return &notification[ddd.AggregateEvent]{
		client: client,
	}
}

func (n notification[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.OrderCreatedEvent:
		return n.onOrderCreated(ctx, event)
	case domain.OrderReadiedEvent:
		return n.onOrderReadied(ctx, event)
	case domain.OrderCanceledEvent:
		return n.onOrderCanceled(ctx, event)
	}
	return nil
}

func (n notification[T]) onOrderCreated(ctx context.Context, event ddd.AggregateEvent) error {
	orderCreated := event.Payload().(*domain.OrderCreated)
	return n.client.NotifyOrderCreated(ctx, event.AggregateID(), orderCreated.CustomerID)
}

func (n notification[T]) onOrderReadied(ctx context.Context, event ddd.AggregateEvent) error {
	orderReadied := event.Payload().(*domain.OrderReadied)
	return n.client.NotifyOrderReady(ctx, event.AggregateID(), orderReadied.CustomerID)
}

func (n notification[T]) onOrderCanceled(ctx context.Context, event ddd.AggregateEvent) error {
	orderCanceled := event.Payload().(*domain.OrderCanceled)
	return n.client.NotifyOrderCanceled(ctx, event.AggregateID(), orderCanceled.CustomerID)
}
