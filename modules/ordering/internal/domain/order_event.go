package domain

import (
	"eda-in-golang/internal/registry"
)

const (
	OrderCreatedEvent   = "ordering.OrderCreated"
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

var _ registry.Registrable = (*OrderCreated)(nil)

func (OrderCreated) Key() string { return OrderCreatedEvent }

type OrderCanceled struct {
	CustomerID string
	PaymentID  string
}

var _ registry.Registrable = (*OrderCanceled)(nil)

func (OrderCanceled) Key() string { return OrderCanceledEvent }

type OrderReadied struct {
	CustomerID string
	PaymentID  string
	Total      float64
}

var _ registry.Registrable = (*OrderReadied)(nil)

func (OrderReadied) Key() string { return OrderReadiedEvent }

type OrderCompleted struct {
	CustomerID string
	InvoiceID  string
}

var _ registry.Registrable = (*OrderCompleted)(nil)

func (OrderCompleted) Key() string { return OrderCompletedEvent }
