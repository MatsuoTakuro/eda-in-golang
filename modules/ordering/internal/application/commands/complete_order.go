package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type CompleteOrderCommander struct {
	orderRepo       infra.OrderRepository
	domainPublisher ddd.EventPublisher
}

func NewCompleteOrderCommander(orderRepo infra.OrderRepository, domainPublisher ddd.EventPublisher) CompleteOrderCommander {
	return CompleteOrderCommander{
		orderRepo:       orderRepo,
		domainPublisher: domainPublisher,
	}
}

func (c CompleteOrderCommander) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	order, err := c.orderRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = order.Complete(cmd.InvoiceID)
	if err != nil {
		return nil
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
