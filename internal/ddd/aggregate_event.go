package ddd

type AggregateEvent interface {
	Event
	AggregateNamer
	AggregateID() string
	AggregateVersion() int
}

type aggregateEvent struct {
	event
}

var _ AggregateEvent = (*aggregateEvent)(nil)

func newAggregateEvent(name string, payload EventPayload, options ...EventOption) AggregateEvent {
	return aggregateEvent{
		event: newEvent(name, payload, options...),
	}
}

func (e aggregateEvent) AggregateName() string { return e.metadata.Get(AggregateNameKey).(string) }
func (e aggregateEvent) AggregateID() string   { return e.metadata.Get(AggregateIDKey).(string) }
func (e aggregateEvent) AggregateVersion() int { return e.metadata.Get(AggregateVersionKey).(int) }
