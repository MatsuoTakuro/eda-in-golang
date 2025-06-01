package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type ApproveOrder struct {
	ID         string
	ShoppingID string
}

type ApproveOrderCommander struct {
	orderRepo infra.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewApproveOrderCommander(orderRepo infra.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) ApproveOrderCommander {
	return ApproveOrderCommander{
		orderRepo: orderRepo,
		publisher: publisher,
	}
}

func (h ApproveOrderCommander) ApproveOrder(ctx context.Context, cmd ApproveOrder) error {
	order, err := h.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Approve(cmd.ShoppingID)
	if err != nil {
		return err
	}

	if err = h.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
