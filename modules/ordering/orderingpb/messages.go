package orderingpb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	OrderAggregateChannel = "mallbots.ordering.events.Order"

	OrderCreatedEvent   = "ordersapi.OrderCreated"
	OrderRejectedEvent  = "ordersapi.OrderRejected"
	OrderApprovedEvent  = "ordersapi.OrderApproved"
	OrderReadiedEvent   = "ordersapi.OrderReadied"
	OrderCanceledEvent  = "ordersapi.OrderCanceled"
	OrderCompletedEvent = "ordersapi.OrderCompleted"

	CommandChannel = "mallbots.ordering.commands"

	RejectOrderCommand  = "ordersapi.RejectOrder"
	ApproveOrderCommand = "ordersapi.ApproveOrder"
)

func Registrations(reg registry.Registry) (err error) {
	regtr := registrar.NewProtoRegistrar(reg)

	// Order events
	if err = regtr.Register(&OrderCreated{}); err != nil {
		return err
	}
	if err = regtr.Register(&OrderRejected{}); err != nil {
		return err
	}
	if err = regtr.Register(&OrderApproved{}); err != nil {
		return err
	}
	if err = regtr.Register(&OrderReadied{}); err != nil {
		return err
	}
	if err = regtr.Register(&OrderCanceled{}); err != nil {
		return err
	}
	if err = regtr.Register(&OrderCompleted{}); err != nil {
		return err
	}

	if err = regtr.Register(&RejectOrder{}); err != nil {
		return err
	}
	if err = regtr.Register(&ApproveOrder{}); err != nil {
		return err
	}

	return nil
}

func (*OrderCreated) Key() string   { return OrderCreatedEvent }
func (*OrderRejected) Key() string  { return OrderRejectedEvent }
func (*OrderApproved) Key() string  { return OrderApprovedEvent }
func (*OrderReadied) Key() string   { return OrderReadiedEvent }
func (*OrderCanceled) Key() string  { return OrderCanceledEvent }
func (*OrderCompleted) Key() string { return OrderCompletedEvent }

func (*RejectOrder) Key() string  { return RejectOrderCommand }
func (*ApproveOrder) Key() string { return ApproveOrderCommand }
