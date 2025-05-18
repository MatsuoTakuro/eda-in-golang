package es

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type EventSourcedAggregate interface {
	ddd.IDer
	AggregateName() string
	ddd.Eventer
	Versioner
	EventApplier
	EventCommitter
}

// AggregateStore rebuilds the state of an aggregate from its events and snapshots
// or saves the events and snapshots of an aggregate to the store.
type AggregateStore interface {
	// Load loads the aggregate from the store and applies the events or snapshot to it.
	Load(ctx context.Context, aggregate EventSourcedAggregate) error
	// Save saves the events of the aggregate to the store and optionally saves a snapshot and publishes the events.
	Save(ctx context.Context, aggregate EventSourcedAggregate) error
}

type AggregateStoreMiddleware func(store AggregateStore) AggregateStore

// AggregateStoreWithMiddleware applies the middleware to the store in reverse order.
// The first middleware in the slice is the outermost, meaning it is the first to enter and the last to exit.
func AggregateStoreWithMiddleware(store AggregateStore, mws ...AggregateStoreMiddleware) AggregateStore {
	//	var s AggregateStore
	s := store
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}
	return s
}
