package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type RejectOrder struct {
	ID string
}

type RejectOrderCommander struct {
	orderRepo infra.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRejectOrderCommander(
	orderRepo infra.OrderRepository,
	publisher ddd.EventPublisher[ddd.Event],
) RejectOrderCommander {
	return RejectOrderCommander{
		orderRepo: orderRepo,
		publisher: publisher,
	}
}

func (c RejectOrderCommander) RejectOrder(ctx context.Context, cmd RejectOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Reject()
	if err != nil {
		return err
	}

	if err = c.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	return c.publisher.Publish(ctx, event)
}
