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
	orderRepo       infra.OrderRepository
	shoppingClient  infra.ShoppingClient
	domainPublisher ddd.EventPublisher
}

func NewCancelOrderCommander(
	orderRepo infra.OrderRepository,
	shoppingClient infra.ShoppingClient,
	domainPublisher ddd.EventPublisher,
) CancelOrderCommander {
	return CancelOrderCommander{
		orderRepo:       orderRepo,
		shoppingClient:  shoppingClient,
		domainPublisher: domainPublisher,
	}
}

func (c CancelOrderCommander) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := c.orderRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Cancel(); err != nil {
		return err
	}

	if err = c.shoppingClient.Cancel(ctx, order.ShoppingID); err != nil {
		return err
	}

	if err = c.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// publish domain events
	if err = c.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	return nil
}
