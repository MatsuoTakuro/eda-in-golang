package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/storespb"
)

func RegisterProductHandlers(productHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	productEvtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return productHandlers.HandleEvent(ctx, eventMsg)
	})

	return stream.Subscribe(storespb.ProductAggregateChannel, productEvtMsgHandler,
		am.MessageFilter{
			storespb.ProductAddedEvent,
			storespb.ProductRebrandedEvent,
			storespb.ProductPriceIncreasedEvent,
			storespb.ProductPriceDecreasedEvent,
			storespb.ProductRemovedEvent,
		},
		am.GroupName("baskets_products_handler"),
	)
}
