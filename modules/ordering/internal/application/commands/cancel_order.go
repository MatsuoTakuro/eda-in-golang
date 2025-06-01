package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type CancelOrder struct {
	ID string
}

type CancelOrderCommander struct {
	orderRepo      infra.OrderRepository
	shoppingClient infra.ShoppingClient
	publisher      ddd.EventPublisher[ddd.Event]
}

func NewCancelOrderCommander(
	orderRepo infra.OrderRepository,
	shoppingClient infra.ShoppingClient,
	publisher ddd.EventPublisher[ddd.Event],
) CancelOrderCommander {
	return CancelOrderCommander{
		orderRepo:      orderRepo,
		shoppingClient: shoppingClient,
		publisher:      publisher,
	}
}

func (c CancelOrderCommander) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Cancel()
	if err != nil {
		return err
	}

	if err = c.shoppingClient.Cancel(ctx, order.ShoppingID); err != nil {
		return err
	}

	if err = c.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	return c.publisher.Publish(ctx, event)
}
