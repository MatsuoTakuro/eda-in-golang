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
	Items      []*domain.Item
}

type CreateOrderCommander struct {
	orderRepo       infra.OrderRepository
	customerClient  infra.CustomerClient
	paymentClient   infra.PaymentClient
	shoppingClient  infra.ShoppingClient
	domainPublisher ddd.EventPublisher
}

func NewCreateOrderCommander(orderRepo infra.OrderRepository,
	customerClient infra.CustomerClient,
	paymentClient infra.PaymentClient,
	shoppingClient infra.ShoppingClient,
	domainPublisher ddd.EventPublisher,
) CreateOrderCommander {
	return CreateOrderCommander{
		orderRepo:       orderRepo,
		customerClient:  customerClient,
		paymentClient:   paymentClient,
		shoppingClient:  shoppingClient,
		domainPublisher: domainPublisher,
	}
}

func (c CreateOrderCommander) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	order, err := domain.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return errors.Wrap(err, "create order command")
	}

	// authorizeCustomer
	if err = c.customerClient.Authorize(ctx, order.CustomerID); err != nil {
		return errors.Wrap(err, "order customer authorization")
	}

	// validatePayment
	if err = c.paymentClient.Confirm(ctx, order.PaymentID); err != nil {
		return errors.Wrap(err, "order payment confirmation")
	}

	// scheduleShopping
	if order.ShoppingID, err = c.shoppingClient.Create(ctx, order); err != nil {
		return errors.Wrap(err, "order shopping scheduling")
	}

	// orderCreation
	if err = c.orderRepo.Save(ctx, order); err != nil {
		return errors.Wrap(err, "order creation")
	}

	// publish domain events
	if err = c.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	return nil
}
