package domain

import (
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
)

type OrderV1 struct {
	CustomerID string
	PaymentID  string
	InvoiceID  string
	ShoppingID string
	Items      []Item
	Status     OrderStatus
}

var _ es.Snapshot = (*OrderV1)(nil)
var _ registry.Registrable = (*OrderV1)(nil)

func (OrderV1) SnapshotName() string { return "ordering.OrderV1" }

func (o OrderV1) Key() string { return o.SnapshotName() }
