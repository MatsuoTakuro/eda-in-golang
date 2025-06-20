package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/baskets/basketspb"
	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/ordering/internal/application"
	"eda-in-golang/modules/ordering/internal/application/commands"
	"eda-in-golang/modules/ordering/internal/domain"

	"github.com/google/uuid"
)

type integrationHandlers[T ddd.Event] struct {
	app application.App
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(app application.App) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		app: app,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err = subscriber.Subscribe(basketspb.BasketAggregateChannel, evtMsgHandler, am.MessageFilter{
		basketspb.BasketCheckedOutEvent,
	}, am.GroupName("ordering-baskets"))
	if err != nil {
		return err
	}

	err = subscriber.Subscribe(depotpb.ShoppingListAggregateChannel, evtMsgHandler, am.MessageFilter{
		depotpb.ShoppingListCompletedEvent,
	}, am.GroupName("ordering-depot"))

	return
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case basketspb.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	case depotpb.ShoppingListCompletedEvent:
		return h.onShoppingListCompleted(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onBasketCheckedOut(ctx context.Context, event ddd.Event) error {
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

	_, _, err := h.app.CreateOrder(ctx, commands.CreateOrder{
		IdempotencyKey: uuid.New().String(), // auto-generated for integration events
		ID:             payload.GetId(),
		CustomerID:     payload.GetCustomerId(),
		PaymentID:      payload.GetPaymentId(),
		Items:          items,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) onShoppingListCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*depotpb.ShoppingListCompleted)

	return h.app.ReadyOrder(ctx, commands.ReadyOrder{ID: payload.GetOrderId()})
}
