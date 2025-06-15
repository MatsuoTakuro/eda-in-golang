package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
)

type CreateOrder struct {
	IdempotencyKey string        `json:"idempotency_key"`
	ID             string        `json:"id"`
	CustomerID     string        `json:"customer_id"`
	PaymentID      string        `json:"payment_id"`
	Items          []domain.Item `json:"items"`
}

type CreateOrderHandler struct {
	events    domain.OrderRepository
	requests  domain.OrderRequestRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateOrderHandler(
	orders domain.OrderRepository,
	orderRequests domain.OrderRequestRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateOrderHandler {
	return CreateOrderHandler{
		events:    orders,
		requests:  orderRequests,
		publisher: publisher,
	}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrder) (
	id string, alreadyAccepted bool, err error,
) {

	if cmd.IdempotencyKey == "" {
		return "", false, errors.Wrap(errors.ErrBadRequest, "idempotency key is required")
	}

	id, inserted, err := h.requests.FindOrInsert(ctx,
		cmd.IdempotencyKey,
		domain.CreateOrderRequest,
		cmd,
	)
	alreadyAccepted = !inserted
	if err != nil {
		return "", alreadyAccepted, err
	}
	if !inserted {
		// If the order request already exists, we can return the existing order ID.
		order, err := h.events.Load(ctx, id)
		if err != nil {
			return "", alreadyAccepted, err
		}
		if order.ID() == "" {
			return "", alreadyAccepted, errors.Wrap(errors.ErrInternal, "order not found")
		}
		if order.ID() != id {
			return "", alreadyAccepted, errors.Wrapf(errors.ErrInternal, "order id mismatch: expected %s, got %s", id, order.ID())
		}

		return order.ID(), alreadyAccepted, nil
	}

	order, err := h.events.Load(ctx, id)
	if err != nil {
		return "", alreadyAccepted, err
	}

	event, err := order.CreateOrder(id, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return "", alreadyAccepted, errors.Wrap(err, "create order command")
	}

	if err = h.events.Save(ctx, order); err != nil {
		return "", alreadyAccepted, errors.Wrap(err, "order creation")
	}

	return id, alreadyAccepted, h.publisher.Publish(ctx, event)
}
