package msghandlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/baskets/basketspb"
	"eda-in-golang/modules/ordering/internal/application"
	"eda-in-golang/modules/ordering/internal/application/commands"
	"eda-in-golang/modules/ordering/internal/domain"
)

func SubscribeIntegrationEvents(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err = subscriber.Subscribe(
		basketspb.BasketAggregateChannel,
		evtMsgHandler,
		am.MessageFilter{
			basketspb.BasketCheckedOutEvent,
		},
		am.GroupName("ordering-baskets"),
	)
	if err != nil {
		return err
	}

	return
}

type integrationHandler[T ddd.Event] struct {
	us application.Usecases
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandler[ddd.Event])(nil)

func NewIntegrationEvents(us application.Usecases) integrationHandler[ddd.Event] {
	return integrationHandler[ddd.Event]{
		us: us,
	}
}

func (h integrationHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case basketspb.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	}

	return nil
}

func (h integrationHandler[T]) onBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*basketspb.BasketCheckedOut)

	items := make([]domain.Item, len(payload.GetItems()))
	for i, item := range payload.GetItems() {
		items[i] = domain.Item{
			ProductID:   item.GetProductId(),
			StoreID:     item.GetStoreId(),
			StoreName:   item.GetStoreName(),
			ProductName: item.GetProductName(),
			Price:       item.GetPrice(),
			Quantity:    int(item.GetQuantity()),
		}
	}

	return h.us.CreateOrder(ctx, commands.CreateOrder{
		ID:         payload.GetId(),
		CustomerID: payload.GetCustomerId(),
		PaymentID:  payload.GetPaymentId(),
		Items:      items,
	})
}
