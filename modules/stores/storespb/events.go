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

// RegisterMessages registers store and product events with the registry.
func RegisterMessages(reg registry.Registry) error {
	regtr := registrar.NewProtoRegistrar(reg)

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

var (
	_ registry.Registrable = (*StoreCreated)(nil)
	_ registry.Registrable = (*StoreParticipationToggled)(nil)
	_ registry.Registrable = (*StoreRebranded)(nil)
	_ registry.Registrable = (*ProductAdded)(nil)
	_ registry.Registrable = (*ProductRebranded)(nil)
	_ registry.Registrable = (*ProductRemoved)(nil)
)

func (*StoreCreated) Key() string              { return StoreCreatedEvent }
func (*StoreParticipationToggled) Key() string { return StoreParticipatingToggledEvent }
func (*StoreRebranded) Key() string            { return StoreRebrandedEvent }
func (*ProductAdded) Key() string              { return ProductAddedEvent }
func (*ProductRebranded) Key() string          { return ProductRebrandedEvent }
func (*ProductRemoved) Key() string            { return ProductRemovedEvent }
