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
	orderRepo       infra.OrderRepository
	domainPublisher ddd.EventPublisher
}

func NewReadyOrderCommander(orderRepo infra.OrderRepository, domainPublisher ddd.EventPublisher) ReadyOrderCommander {
	return ReadyOrderCommander{
		orderRepo:       orderRepo,
		domainPublisher: domainPublisher,
	}
}

func (c ReadyOrderCommander) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := c.orderRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Ready(); err != nil {
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
