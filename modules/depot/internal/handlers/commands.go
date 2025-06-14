package handlers

import (
	"context"
	"database/sql"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/depot/depotpb"
)

func SubscribeCommands(container di.Container) error {
	cmdMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.AckableRawMessage) (err error) {
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

		cmdMsgHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewCommandMessageHandler(
				di.Get(ctx, di.Registry).(registry.Registry),
				di.Get(ctx, di.ReplyStream).(am.ReplyStream),
				di.Get(ctx, di.CommandHandler).(ddd.CommandHandler[ddd.Command]),
			),
			di.Get(ctx, di.InboxMiddleware).(am.RawMessageHandlerMiddleware),
		)

		return cmdMsgHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get(di.Stream).(am.RawMessageStream)

	return subscribeCommands(subscriber, am.MessageHandlerFunc[am.AckableRawMessage](cmdMsgHandler))
}

func subscribeCommands(subscriber am.RawMessageSubscriber, handler am.RawMessageHandler) error {
	return subscriber.Subscribe(depotpb.CommandChannel,
		handler,
		am.MessageFilter{
			depotpb.CreateShoppingListCommand,
			depotpb.CancelShoppingListCommand,
			depotpb.InitiateShoppingCommand,
		},
		am.GroupName("depot-commands"),
	)
}
