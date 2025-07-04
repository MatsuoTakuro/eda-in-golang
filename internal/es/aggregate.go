package es

import (
	"eda-in-golang/internal/ddd"
)

type Aggregate interface {
	ddd.Aggregate
	EventCommitter
	Versioner
	VersionSetter
}

type Versioner interface {
	// Version returns the current committed version of the aggregate.
	Version() int
	// PendingVersion returns the version the aggregate will have after uncommitted events are applied.
	PendingVersion() int
}

type aggregate struct {
	ddd.Aggregate
	version int
}

var _ Aggregate = (*aggregate)(nil)

func NewAggregate(id, name string) *aggregate {
	return &aggregate{
		Aggregate: ddd.NewAggregate(id, name),
		version:   0,
	}
}

func (a *aggregate) AddEvent(name string, payload ddd.EventPayload, options ...ddd.EventOption) {
	options = append(
		options,
		ddd.Metadata{
			ddd.AggregateVersionKey: a.PendingVersion() + 1,
		},
	)
	a.Aggregate.AddEvent(name, payload, options...)
}

func (a *aggregate) CommitEvents() {
	a.version += len(a.Events())
	a.ClearEvents()
}

func (a *aggregate) Version() int           { return a.version }
func (a *aggregate) PendingVersion() int    { return a.version + len(a.Events()) }
func (a *aggregate) SetVersion(version int) { a.version = version }
