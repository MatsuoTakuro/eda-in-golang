package msghandlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/orderingpb"
)

func SubscribeDomainEvents(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.OrderCreatedEvent,
		domain.OrderRejectedEvent,
		domain.OrderApprovedEvent,
		domain.OrderReadiedEvent,
		domain.OrderCanceledEvent,
		domain.OrderCompletedEvent,
	)
}

type domainEvtHdlr[T ddd.Event] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.Event] = (*domainEvtHdlr[ddd.Event])(nil)

func NewDomainEvents(publisher am.MessagePublisher[ddd.Event]) domainEvtHdlr[ddd.Event] {
	return domainEvtHdlr[ddd.Event]{publisher: publisher}
}

func (h domainEvtHdlr[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case domain.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case domain.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	case domain.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)
	}
	return nil
}

func (h domainEvtHdlr[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)
	items := make([]*orderingpb.OrderCreated_Item, len(payload.Items))
	for i, item := range payload.Items {
		items[i] = &orderingpb.OrderCreated_Item{
			ProductId: item.ProductID,
			StoreId:   item.StoreID,
			Price:     item.Price,
			Quantity:  int32(item.Quantity),
		}
	}
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderCreatedEvent, &orderingpb.OrderCreated{
			Id:         payload.ID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
			ShoppingId: payload.ShoppingID,
			Items:      items,
		}),
	)
}

func (h domainEvtHdlr[T]) onOrderRejected(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderRejectedEvent, &orderingpb.OrderRejected{
			Id:         payload.ID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
		}),
	)
}

func (h domainEvtHdlr[T]) onOrderApproved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderApprovedEvent, &orderingpb.OrderApproved{
			Id:         payload.ID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
		}),
	)
}

func (h domainEvtHdlr[T]) onOrderReadied(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderReadiedEvent, &orderingpb.OrderReadied{
			Id:         payload.ID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
			Total:      payload.GetTotal(),
		}),
	)
}

func (h domainEvtHdlr[T]) onOrderCanceled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderCanceledEvent, &orderingpb.OrderCanceled{
			Id:         payload.ID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
		}),
	)
}

func (h domainEvtHdlr[T]) onOrderCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderCompletedEvent, &orderingpb.OrderCompleted{
			Id:         payload.ID(),
			CustomerId: payload.CustomerID,
			InvoiceId:  payload.InvoiceID,
		}),
	)
}
