package application

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/customers/customerspb"
)

type commandHandlers struct {
	app App
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandlers)(nil)

func NewCommandHandlers(app App) commandHandlers {
	return commandHandlers{
		app: app,
	}
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case customerspb.AuthorizeCustomerCommand:
		return h.doAuthorizeCustomer(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doAuthorizeCustomer(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*customerspb.AuthorizeCustomer)

	// no reply
	return nil, h.app.AuthorizeCustomer(ctx, AuthorizeCustomer{ID: payload.GetId()})
}
