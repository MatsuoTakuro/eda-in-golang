package es

import (
	"fmt"

	"eda-in-golang/internal/ddd"
)

// Hydrator applies events and snapshots to rebuild the state of an aggregate.
type Hydrator interface {
	EventApplier
	Snapshotter
}

type EventApplier interface {
	ApplyEvent(ddd.Event) error
}

type EventCommitter interface {
	// CommitEvents updates the version of the aggregate and clears the events.
	// It must be called after all events that were added to the aggregate have been saved to the store.
	CommitEvents()
}

// LoadEvent applies an event to an aggregate. It is used to rebuild the state of an aggregate from its events.
func LoadEvent(v interface{}, event ddd.AggregateEvent) error {
	type loader interface {
		EventApplier
		VersionSetter
	}

	agg, ok := v.(loader)
	if !ok {
		return fmt.Errorf("%T does not have the methods implemented to load events", v)
	}

	if err := agg.ApplyEvent(event); err != nil {
		return err
	}
	agg.SetVersion(event.AggregateVersion())

	return nil
}
