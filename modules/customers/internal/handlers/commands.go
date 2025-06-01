package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/customers/customerspb"
)

func SubscribeCommands(subscriber am.CommandSubscriber, handlers ddd.CommandHandler[ddd.Command]) error {
	cmdMsgHandler := am.CommandMessageHandlerFunc(func(ctx context.Context, cmdMsg am.CommandMessage) (ddd.Reply, error) {
		return handlers.HandleCommand(ctx, cmdMsg)
	})

	return subscriber.Subscribe(customerspb.CommandChannel, cmdMsgHandler,
		am.MessageFilter{
			customerspb.AuthorizeCustomerCommand,
		},
		am.GroupName("customer-commands"),
	)
}
