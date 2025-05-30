package ddd

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	IDer
	EventName() string
	Payload() EventPayload
	Metadata() Metadata
	OccurredAt() time.Time
}

type EventOption interface {
	configureEvent(*event)
}

type EventPayload any

func NewEvent(name string, payload EventPayload, options ...EventOption) event {
	return newEvent(name, payload, options...)
}

type event struct {
	Entity
	payload    EventPayload
	metadata   Metadata
	occurredAt time.Time
}

var _ Event = (*event)(nil)

func newEvent(name string, payload EventPayload, options ...EventOption) event {
	evt := event{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, opt := range options {
		opt.configureEvent(&evt)
	}

	return evt
}

func (e event) EventName() string     { return e.EntityName() }
func (e event) Payload() EventPayload { return e.payload }
func (e event) Metadata() Metadata    { return e.metadata }
func (e event) OccurredAt() time.Time { return e.occurredAt }
