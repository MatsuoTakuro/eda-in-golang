package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/storespb"
)

func RegisterStoreHandler(storeHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	storeEvtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return storeHandlers.HandleEvent(ctx, eventMsg)
	})

	return stream.Subscribe(storespb.StoreAggregateChannel, storeEvtMsgHandler,
		am.MessageFilter{
			storespb.StoreCreatedEvent,
			storespb.StoreParticipatingToggledEvent,
			storespb.StoreRebrandedEvent,
		},
		am.GroupName("baskets_stores_handler"),
	)
}
