package commands

import (
	"context"

	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type CancelOrder struct {
	ID string
}

type CancelOrderCommander struct {
	orderRepo      infra.OrderRepository
	shoppingClient infra.ShoppingClient
}

func NewCancelOrderCommander(
	orderRepo infra.OrderRepository,
	shoppingClient infra.ShoppingClient,
) CancelOrderCommander {
	return CancelOrderCommander{
		orderRepo:      orderRepo,
		shoppingClient: shoppingClient,
	}
}

func (c CancelOrderCommander) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Cancel(); err != nil {
		return err
	}

	if err = c.shoppingClient.Cancel(ctx, order.ShoppingID); err != nil {
		return err
	}

	if err = c.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
