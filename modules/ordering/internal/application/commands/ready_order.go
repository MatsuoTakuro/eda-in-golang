package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type ReadyOrder struct {
	ID string
}

type ReadyOrderCommander struct {
	orderRepo infra.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewReadyOrderCommander(
	orderRepo infra.OrderRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ReadyOrderCommander {
	return ReadyOrderCommander{
		orderRepo: orderRepo,
		publisher: publisher,
	}
}

func (c ReadyOrderCommander) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Ready()
	if err != nil {
		return nil
	}

	if err = c.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	return c.publisher.Publish(ctx, event)
}
