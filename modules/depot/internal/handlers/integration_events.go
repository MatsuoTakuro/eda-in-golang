package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/storespb"
)

func SubscribeIntegrationEvents(
	subscriber am.EventSubscriber,
	handlers ddd.EventHandler[ddd.Event],
) (err error) {

	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](
		func(ctx context.Context, eventMsg am.EventMessage) error {
			return handlers.HandleEvent(ctx, eventMsg)
		},
	)

	err = subscriber.Subscribe(
		storespb.StoreAggregateChannel,
		evtMsgHandler,
		am.MessageFilter{
			storespb.StoreCreatedEvent,
			storespb.StoreRebrandedEvent,
		},
		am.GroupName("depot-stores"),
	)
	if err != nil {
		return err
	}

	err = subscriber.Subscribe(
		storespb.ProductAggregateChannel,
		evtMsgHandler,
		am.MessageFilter{
			storespb.ProductAddedEvent,
			storespb.ProductRebrandedEvent,
			storespb.ProductPriceIncreasedEvent,
			storespb.ProductPriceDecreasedEvent,
			storespb.ProductRemovedEvent,
		},
		am.GroupName("depot-products"),
	)
	if err != nil {
		return err
	}

	return nil
}
