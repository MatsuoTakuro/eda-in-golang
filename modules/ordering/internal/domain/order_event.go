package domain

import (
	"eda-in-golang/internal/registry"
)

const (
	OrderCreatedEvent   = "ordering.OrderCreated"
	OrderRejectedEvent  = "ordering.OrderRejected"
	OrderApprovedEvent  = "ordering.OrderApproved"
	OrderCanceledEvent  = "ordering.OrderCanceled"
	OrderReadiedEvent   = "ordering.OrderReadied"
	OrderCompletedEvent = "ordering.OrderCompleted"
)

type OrderCreated struct {
	CustomerID string
	PaymentID  string
	ShoppingID string
	Items      []Item
}

type OrderRejected struct{}

type OrderApproved struct {
	ShoppingID string
}

type OrderCanceled struct {
	CustomerID string
	PaymentID  string
}

type OrderReadied struct {
	CustomerID string
	PaymentID  string
	Total      float64
}

type OrderCompleted struct {
	CustomerID string
	InvoiceID  string
}

var (
	_ registry.Registrable = (*OrderCreated)(nil)
	_ registry.Registrable = (*OrderRejected)(nil)
	_ registry.Registrable = (*OrderApproved)(nil)
	_ registry.Registrable = (*OrderCanceled)(nil)
	_ registry.Registrable = (*OrderReadied)(nil)
	_ registry.Registrable = (*OrderCompleted)(nil)
)

func (OrderCreated) Key() string   { return OrderCreatedEvent }
func (OrderRejected) Key() string  { return OrderRejectedEvent }
func (OrderApproved) Key() string  { return OrderApprovedEvent }
func (OrderCanceled) Key() string  { return OrderCanceledEvent }
func (OrderReadied) Key() string   { return OrderReadiedEvent }
func (OrderCompleted) Key() string { return OrderCompletedEvent }
