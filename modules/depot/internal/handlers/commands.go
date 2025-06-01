package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/depotpb"
)

func SubscribeCommands(subscriber am.CommandSubscriber, handlers ddd.CommandHandler[ddd.Command]) error {
	cmdMsgHandler := am.CommandMessageHandlerFunc(func(ctx context.Context, cmdMsg am.CommandMessage) (ddd.Reply, error) {
		return handlers.HandleCommand(ctx, cmdMsg)
	})

	return subscriber.Subscribe(depotpb.CommandChannel, cmdMsgHandler,
		am.MessageFilter{
			depotpb.CreateShoppingListCommand,
			depotpb.CancelShoppingListCommand,
			depotpb.InitiateShoppingCommand,
		},
		am.GroupName("depot-commands"),
	)
}
