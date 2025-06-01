package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type CreateOrder struct {
	ID         string
	CustomerID string
	PaymentID  string
	Items      []domain.Item
}

type CreateOrderCommander struct {
	orderRepo infra.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateOrderCommander(
	orderRepo infra.OrderRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateOrderCommander {
	return CreateOrderCommander{
		orderRepo: orderRepo,
		publisher: publisher,
	}
}

func (c CreateOrderCommander) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return errors.Wrap(err, "create order command")
	}

	if err = c.orderRepo.Save(ctx, order); err != nil {
		return errors.Wrap(err, "order creation")
	}

	return c.publisher.Publish(ctx, event)
}
