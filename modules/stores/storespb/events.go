package storespb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	StoreAggregateChannel = "mallbots.stores.events.Store"

	StoreCreatedEvent              = "storesapi.StoreCreated"
	StoreParticipatingToggledEvent = "storesapi.StoreParticipatingToggled"
	StoreRebrandedEvent            = "storesapi.StoreRebranded"

	ProductAggregateChannel = "mallbots.stores.events.Product"

	ProductAddedEvent          = "storesapi.ProductAdded"
	ProductRebrandedEvent      = "storesapi.ProductRebranded"
	ProductPriceIncreasedEvent = "storesapi.ProductPriceIncreased"
	ProductPriceDecreasedEvent = "storesapi.ProductPriceDecreased"
	ProductRemovedEvent        = "storesapi.ProductRemoved"
)

func RegisterMessages(reg registry.Registry) error {
	return RegisterMessagesWithRegistrar(registrar.NewProtoRegistrar(reg))
}

func RegisterMessagesWithRegistrar(regtr registry.Registrar) error {
	// Store events
	if err := regtr.Register(&StoreCreated{}); err != nil {
		return err
	}
	if err := regtr.Register(&StoreParticipationToggled{}); err != nil {
		return err
	}
	if err := regtr.Register(&StoreRebranded{}); err != nil {
		return err
	}

	if err := regtr.Register(&ProductAdded{}); err != nil {
		return err
	}
	if err := regtr.Register(&ProductRebranded{}); err != nil {
		return err
	}
	if err := regtr.RegisterWithKey(ProductPriceIncreasedEvent, &ProductPriceChanged{}); err != nil {
		return err
	}
	if err := regtr.RegisterWithKey(ProductPriceDecreasedEvent, &ProductPriceChanged{}); err != nil {
		return err
	}
	if err := regtr.Register(&ProductRemoved{}); err != nil {
		return err
	}

	return nil
}

func (*StoreCreated) Key() string              { return StoreCreatedEvent }
func (*StoreParticipationToggled) Key() string { return StoreParticipatingToggledEvent }
func (*StoreRebranded) Key() string            { return StoreRebrandedEvent }

func (*ProductAdded) Key() string     { return ProductAddedEvent }
func (*ProductRebranded) Key() string { return ProductRebrandedEvent }
func (*ProductRemoved) Key() string   { return ProductRemovedEvent }
