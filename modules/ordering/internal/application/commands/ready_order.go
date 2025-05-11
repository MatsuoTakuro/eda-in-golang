package commands

import (
	"context"

	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type ReadyOrder struct {
	ID string
}

type ReadyOrderCommander struct {
	orderRepo infra.OrderRepository
}

func NewReadyOrderCommander(orderRepo infra.OrderRepository) ReadyOrderCommander {
	return ReadyOrderCommander{
		orderRepo: orderRepo,
	}
}

func (c ReadyOrderCommander) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Ready(); err != nil {
		return nil
	}

	if err = c.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
