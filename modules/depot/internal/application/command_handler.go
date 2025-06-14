package application

import (
	"context"

	"github.com/google/uuid"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/depot/internal/application/commands"
)

type commandHandler struct {
	app App
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandler)(nil)

// NewCommandHandler creates a handler that handle commands mainly coming from other modules.
func NewCommandHandler(app App) commandHandler {
	return commandHandler{
		app: app,
	}
}

func (h commandHandler) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case depotpb.CreateShoppingListCommand:
		return h.doCreateShoppingList(ctx, cmd)
	case depotpb.CancelShoppingListCommand:
		return h.doCancelShoppingList(ctx, cmd)
	case depotpb.InitiateShoppingCommand:
		return h.doInitiateShopping(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandler) doCreateShoppingList(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*depotpb.CreateShoppingList)

	id := uuid.New().String()

	items := make([]commands.OrderItem, 0, len(payload.GetItems()))
	for _, item := range payload.GetItems() {
		items = append(items, commands.OrderItem{
			StoreID:   item.GetStoreId(),
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
		})
	}

	err := h.app.CreateShoppingList(ctx, commands.CreateShoppingList{
		ID:      id,
		OrderID: payload.GetOrderId(),
		Items:   items,
	})

	return ddd.NewReply(depotpb.CreatedShoppingListReply, &depotpb.CreatedShoppingList{Id: id}), err
}

func (h commandHandler) doCancelShoppingList(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*depotpb.CancelShoppingList)

	err := h.app.CancelShoppingList(ctx, commands.CancelShoppingList{ID: payload.GetId()})

	// returning nil returns a simple Success or Failure reply; err being nil determines which
	return nil, err
}

func (h commandHandler) doInitiateShopping(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*depotpb.InitiateShopping)

	err := h.app.InitiateShopping(ctx, commands.InitiateShopping{ID: payload.GetId()})

	// returning nil returns a simple Success or Failure reply; err being nil determines which
	return nil, err
}
