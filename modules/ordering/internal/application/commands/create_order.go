package commands

import (
	"context"

	"github.com/stackus/errors"

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
	orderRepo   infra.OrderRepository
	customerCli infra.CustomerClient
	paymentCli  infra.PaymentClient
	shoppingCli infra.ShoppingClient
}

func NewCreateOrderCommander(orderRepo infra.OrderRepository,
	customerCli infra.CustomerClient,
	paymentCli infra.PaymentClient,
	shoppingCli infra.ShoppingClient,
) CreateOrderCommander {
	return CreateOrderCommander{
		orderRepo:   orderRepo,
		customerCli: customerCli,
		paymentCli:  paymentCli,
		shoppingCli: shoppingCli,
	}
}

func (c CreateOrderCommander) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	order, err := c.orderRepo.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// authorizeCustomer
	if err = c.customerCli.Authorize(ctx, cmd.CustomerID); err != nil {
		return errors.Wrap(err, "order customer authorization")
	}

	// validatePayment
	if err = c.paymentCli.Confirm(ctx, cmd.PaymentID); err != nil {
		return errors.Wrap(err, "order payment confirmation")
	}

	// scheduleShopping
	var shoppingID string
	if shoppingID, err = c.shoppingCli.Create(ctx, cmd.ID, cmd.Items); err != nil {
		return errors.Wrap(err, "order shopping scheduling")
	}

	err = order.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, shoppingID, cmd.Items)
	if err != nil {
		return errors.Wrap(err, "create order command")
	}

	// orderCreation
	if err = c.orderRepo.Save(ctx, order); err != nil {
		return errors.Wrap(err, "order creation")
	}

	return nil
}
