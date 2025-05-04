package eventhandlers

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type Notification interface {
	OnOrderCreatedEventHandler
	OnOrderReadiedEventHandler
	OnOrderCanceledEventHandler
}

type notification struct {
	client infra.NotificationClient
}

var _ Notification = (*notification)(nil)

func NewNotification(client infra.NotificationClient) Notification {
	return &notification{
		client: client,
	}
}

func (h notification) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	orderCreated := event.(*domain.OrderCreated)
	return h.client.NotifyOrderCreated(ctx, orderCreated.Order.ID, orderCreated.Order.CustomerID)
}

func (h notification) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	orderReadied := event.(*domain.OrderReadied)
	return h.client.NotifyOrderReady(ctx, orderReadied.Order.ID, orderReadied.Order.CustomerID)
}

func (h notification) OnOrderCanceled(ctx context.Context, event ddd.Event) error {
	orderCanceled := event.(*domain.OrderCanceled)
	return h.client.NotifyOrderCanceled(ctx, orderCanceled.Order.ID, orderCanceled.Order.CustomerID)
}
