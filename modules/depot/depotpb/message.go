package depotpb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	ShoppingListAggregateChannel = "mallbots.depot.events.ShoppingList"

	ShoppingListCompletedEvent = "depotapi.ShoppingListCompleted"
)

const (
	CommandChannel = "mallbots.depot.commands"

	CreateShoppingListCommand = "depotapi.CreateShoppingListCommand"
	CancelShoppingListCommand = "depotapi.CancelShoppingListCommand"
	InitiateShoppingCommand   = "depotapi.InitiateShoppingCommand"
)

const (
	CreatedShoppingListReply = "depotapi.CreatedShoppingListReply"
)

func RegisterMessages(reg registry.Registry) (err error) {
	regtr := registrar.NewProtoRegistrar(reg)

	if err = regtr.Register(&ShoppingListCompleted{}); err != nil {
		return
	}

	if err = regtr.Register(&CreateShoppingList{}); err != nil {
		return err
	}
	if err = regtr.Register(&CancelShoppingList{}); err != nil {
		return err
	}
	if err = regtr.Register(&InitiateShopping{}); err != nil {
		return err
	}

	if err = regtr.Register(&CreatedShoppingList{}); err != nil {
		return err
	}

	return nil
}

// Events
var (
	_ registry.Registrable = (*ShoppingListCompleted)(nil)
)

func (*ShoppingListCompleted) Key() string { return ShoppingListCompletedEvent }

// Commands
var (
	_ registry.Registrable = (*CreateShoppingList)(nil)
	_ registry.Registrable = (*CancelShoppingList)(nil)
	_ registry.Registrable = (*InitiateShopping)(nil)
)

func (*CreateShoppingList) Key() string { return CreateShoppingListCommand }
func (*CancelShoppingList) Key() string { return CancelShoppingListCommand }
func (*InitiateShopping) Key() string   { return InitiateShoppingCommand }

// Replies
var (
	_ registry.Registrable = (*CreatedShoppingList)(nil)
)

func (*CreatedShoppingList) Key() string { return CreatedShoppingListReply }
