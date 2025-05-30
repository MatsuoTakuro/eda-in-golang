package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/baskets/internal/domain"
)

type OrderHandlers[T ddd.AggregateEvent] struct {
	orders domain.OrderClient
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*OrderHandlers[ddd.AggregateEvent])(nil)

func NewOrderHandlers(orders domain.OrderClient) OrderHandlers[ddd.AggregateEvent] {
	return OrderHandlers[ddd.AggregateEvent]{
		orders: orders,
	}
}

func (h OrderHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	}
	return nil
}

func (h OrderHandlers[T]) onBasketCheckedOut(ctx context.Context, event ddd.AggregateEvent) error {
	checkedOut := event.Payload().(*domain.BasketCheckedOut)
	_, err := h.orders.Save(ctx, checkedOut.PaymentID, checkedOut.CustomerID, checkedOut.Items)
	return err
}
