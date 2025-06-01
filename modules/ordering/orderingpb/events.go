package orderingpb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	OrderAggregateChannel = "mallbots.ordering.events.Order"

	OrderCreatedEvent   = "ordersapi.OrderCreated"
	OrderReadiedEvent   = "ordersapi.OrderReadied"
	OrderCanceledEvent  = "ordersapi.OrderCanceled"
	OrderCompletedEvent = "ordersapi.OrderCompleted"
)

func RegisterMessages(reg registry.Registry) error {
	regtr := registrar.NewProtoRegistrar(reg)

	// Order events
	if err := regtr.Register(&OrderCreated{}); err != nil {
		return err
	}
	if err := regtr.Register(&OrderReadied{}); err != nil {
		return err
	}
	if err := regtr.Register(&OrderCanceled{}); err != nil {
		return err
	}
	if err := regtr.Register(&OrderCompleted{}); err != nil {
		return err
	}

	return nil
}

var (
	_ registry.Registrable = (*OrderCreated)(nil)
	_ registry.Registrable = (*OrderReadied)(nil)
	_ registry.Registrable = (*OrderCanceled)(nil)
	_ registry.Registrable = (*OrderCompleted)(nil)
)

func (*OrderCreated) Key() string   { return OrderCreatedEvent }
func (*OrderReadied) Key() string   { return OrderReadiedEvent }
func (*OrderCanceled) Key() string  { return OrderCanceledEvent }
func (*OrderCompleted) Key() string { return OrderCompletedEvent }
