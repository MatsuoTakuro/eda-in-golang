package ddd

const (
	AggregateNameKey    = "aggregate-name"
	AggregateIDKey      = "aggregate-id"
	AggregateVersionKey = "aggregate-version"
)

type Aggregate interface {
	Entity
	AggregateNamer
	Eventer
}

type AggregateNamer interface {
	AggregateName() string
}

type Eventer interface {
	AddEvent(string, EventPayload, ...EventOption)
	Events() []AggregateEvent
	ClearEvents()
}

type aggregate struct {
	Entity
	events []AggregateEvent
}

var _ Aggregate = (*aggregate)(nil)

func NewAggregate(id, name string) *aggregate {
	return &aggregate{
		Entity: NewEntity(id, name),
		events: make([]AggregateEvent, 0),
	}
}

func (a *aggregate) AggregateName() string { return a.EntityName() }
func (a *aggregate) AddEvent(name string, payload EventPayload, options ...EventOption) {
	options = append(
		options,
		Metadata{
			AggregateNameKey: a.EntityName(),
			AggregateIDKey:   a.ID(),
		},
	)
	a.events = append(
		a.events,
		newAggregateEvent(name, payload, options...),
	)
}
func (a *aggregate) Events() []AggregateEvent { return a.events }
func (a *aggregate) ClearEvents()             { a.events = []AggregateEvent{} }
