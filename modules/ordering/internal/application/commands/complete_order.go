package commands

import (
	"context"

	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type CompleteOrderCommander struct {
	orderRepo infra.OrderRepository
}

func NewCompleteOrderCommander(orderRepo infra.OrderRepository) CompleteOrderCommander {
	return CompleteOrderCommander{
		orderRepo: orderRepo,
	}
}

func (c CompleteOrderCommander) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = order.Complete(cmd.InvoiceID)
	if err != nil {
		return nil
	}
	if err = c.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
