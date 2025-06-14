package handlers

import (
	"context"
	"database/sql"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/stores/storespb"
)

func SubscribeIntegrationEvents(container di.Container) error {
	evtMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.AckableRawMessage) (err error) {
		ctx = container.Scoped(ctx)
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}(di.Get(ctx, di.TX).(*sql.Tx))

		evtHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewEventMessageHandler(
				di.Get(ctx, di.Registry).(registry.Registry),
				di.Get(ctx, di.IntegrationEventHandler).(ddd.EventHandler[ddd.Event]),
			),
			di.Get(ctx, di.InboxMiddleware).(am.RawMessageHandlerMiddleware),
		)

		return evtHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get(di.Stream).(am.RawMessageStream)

	err := subscriber.Subscribe(storespb.StoreAggregateChannel,
		am.MessageHandlerFunc[am.AckableRawMessage](evtMsgHandler),
		am.MessageFilter{
			storespb.StoreCreatedEvent,
			storespb.StoreRebrandedEvent,
		},
		am.GroupName("depot-stores"),
	)
	if err != nil {
		return err
	}

	err = subscriber.Subscribe(storespb.ProductAggregateChannel,
		am.MessageHandlerFunc[am.AckableRawMessage](evtMsgHandler),
		am.MessageFilter{
			storespb.ProductAddedEvent,
			storespb.ProductRebrandedEvent,
			storespb.ProductRemovedEvent,
		},
		am.GroupName("depot-products"),
	)
	if err != nil {
		return err
	}

	return nil
}
