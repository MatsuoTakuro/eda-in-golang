package msghandlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/application"
	"eda-in-golang/modules/ordering/internal/application/commands"
	"eda-in-golang/modules/ordering/orderingpb"
)

func SubscribeCommands(
	subscriber am.CommandSubscriber,
	handlers ddd.CommandHandler[ddd.Command],
) error {

	cmdMsgHandler := am.CommandMessageHandlerFunc(
		func(ctx context.Context, cmdMsg am.CommandMessage) (ddd.Reply, error) {
			return handlers.HandleCommand(ctx, cmdMsg)
		},
	)

	return subscriber.Subscribe(
		orderingpb.CommandChannel,
		cmdMsgHandler,
		am.MessageFilter{
			orderingpb.RejectOrderCommand,
			orderingpb.ApproveOrderCommand,
		},
		am.GroupName("ordering-commands"),
	)
}

type cmdHdlr struct {
	us application.Usecases
}

var _ ddd.CommandHandler[ddd.Command] = (*cmdHdlr)(nil)

func NewCommands(us application.Usecases) cmdHdlr {
	return cmdHdlr{
		us: us,
	}
}

func (h cmdHdlr) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case orderingpb.RejectOrderCommand:
		return h.doRejectOrder(ctx, cmd)
	case orderingpb.ApproveOrderCommand:
		return h.doApproveOrder(ctx, cmd)
	}

	return nil, nil
}

func (h cmdHdlr) doRejectOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.RejectOrder)

	return nil, h.us.RejectOrder(ctx, commands.RejectOrder{ID: payload.GetId()})
}

func (h cmdHdlr) doApproveOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.ApproveOrder)

	return nil, h.us.ApproveOrder(ctx, commands.ApproveOrder{
		ID:         payload.GetId(),
		ShoppingID: payload.GetShoppingId(),
	})
}
