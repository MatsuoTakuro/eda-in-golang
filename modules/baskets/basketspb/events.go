package basketspb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	BasketAggregateChannel = "mallbots.baskets.events.Basket"

	BasketStartedEvent    = "basketsapi.BasketStarted"
	BasketCanceledEvent   = "basketsapi.BasketCanceled"
	BasketCheckedOutEvent = "basketsapi.BasketCheckedOut"
)

func RegisterMessages(reg registry.Registry) error {
	regtr := registrar.NewProtoRegistrar(reg)

	// Basket events
	if err := regtr.Register(&BasketStarted{}); err != nil {
		return err
	}
	if err := regtr.Register(&BasketCanceled{}); err != nil {
		return err
	}
	if err := regtr.Register(&BasketCheckedOut{}); err != nil {
		return err
	}

	return nil
}

func (*BasketStarted) Key() string    { return BasketStartedEvent }
func (*BasketCanceled) Key() string   { return BasketCanceledEvent }
func (*BasketCheckedOut) Key() string { return BasketCheckedOutEvent }
