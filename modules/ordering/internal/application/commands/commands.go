package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type Commands interface {
	CreateOrder(ctx context.Context, cmd CreateOrder) error
	CancelOrder(ctx context.Context, cmd CancelOrder) error
	ReadyOrder(ctx context.Context, cmd ReadyOrder) error
	CompleteOrder(ctx context.Context, cmd CompleteOrder) error
}

type commands struct {
	CreateOrderCommander
	CancelOrderCommander
	ReadyOrderCommander
	CompleteOrderCommander
}

func New(orderRepo infra.OrderRepository,
	customerClient infra.CustomerClient, paymentClient infra.PaymentClient, shoppingClient infra.ShoppingClient,
	domainPublisher ddd.EventPublisher,
) Commands {
	return &commands{
		CreateOrderCommander:   NewCreateOrderCommander(orderRepo, customerClient, paymentClient, shoppingClient, domainPublisher),
		CancelOrderCommander:   NewCancelOrderCommander(orderRepo, shoppingClient, domainPublisher),
		ReadyOrderCommander:    NewReadyOrderCommander(orderRepo, domainPublisher),
		CompleteOrderCommander: NewCompleteOrderCommander(orderRepo, domainPublisher),
	}
}
