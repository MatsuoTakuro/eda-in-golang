package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type Commands interface {
	CreateOrder(ctx context.Context, cmd CreateOrder) error
	RejectOrder(ctx context.Context, cmd RejectOrder) error
	ApproveOrder(ctx context.Context, cmd ApproveOrder) error
	CancelOrder(ctx context.Context, cmd CancelOrder) error
	ReadyOrder(ctx context.Context, cmd ReadyOrder) error
	CompleteOrder(ctx context.Context, cmd CompleteOrder) error
}
type commands struct {
	CreateOrderCommander
	RejectOrderCommander
	ApproveOrderCommander
	CancelOrderCommander
	ReadyOrderCommander
	CompleteOrderCommander
}

func New(
	orderRepo infra.OrderRepository,
	shoppingClient infra.ShoppingClient,
	publisher ddd.EventPublisher[ddd.Event],
) Commands {
	return &commands{
		CreateOrderCommander:   NewCreateOrderCommander(orderRepo, publisher),
		RejectOrderCommander:   NewRejectOrderCommander(orderRepo, publisher),
		ApproveOrderCommander:  NewApproveOrderCommander(orderRepo, publisher),
		CancelOrderCommander:   NewCancelOrderCommander(orderRepo, shoppingClient, publisher),
		ReadyOrderCommander:    NewReadyOrderCommander(orderRepo, publisher),
		CompleteOrderCommander: NewCompleteOrderCommander(orderRepo, publisher),
	}
}
