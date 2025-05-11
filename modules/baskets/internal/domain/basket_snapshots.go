package domain

import "eda-in-golang/internal/es"

type BasketV1 struct {
	CustomerID string
	PaymentID  string
	Items      map[string]Item
	Status     BasketStatus
}

var _ es.Snapshot = (*BasketV1)(nil)

func (BasketV1) SnapshotName() string { return "baskets.BasketV1" }

func (b BasketV1) Key() string { return b.SnapshotName() }
