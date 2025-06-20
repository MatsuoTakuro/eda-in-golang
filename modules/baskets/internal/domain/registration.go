package domain

import (
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

func Registrations(reg registry.Registry) error {
	regtr := registrar.NewJsonRegistrar(reg)

	// Basket
	if err := regtr.Register(Basket{}, func(v interface{}) error {
		basket := v.(*Basket)
		basket.Aggregate = es.NewAggregate("", BasketAggregate)
		basket.Items = make(map[string]Item)
		return nil
	}); err != nil {
		return err
	}
	// basket events
	if err := regtr.Register(BasketStarted{}); err != nil {
		return err
	}
	if err := regtr.Register(BasketCanceled{}); err != nil {
		return err
	}
	if err := regtr.Register(BasketCheckedOut{}); err != nil {
		return err
	}
	if err := regtr.Register(BasketItemAdded{}); err != nil {
		return err
	}
	if err := regtr.Register(BasketItemRemoved{}); err != nil {
		return err
	}
	// basket snapshots
	if err := regtr.RegisterWithKey(BasketV1{}.SnapshotName(), BasketV1{}); err != nil {
		return err
	}

	return nil
}
