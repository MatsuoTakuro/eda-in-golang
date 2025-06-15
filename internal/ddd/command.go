package ddd

import (
	"time"

	"github.com/google/uuid"
)

type CommandOption interface {
	configureCommand(*command)
}

type CommandPayload any

type Command interface {
	IDer
	CommandName() string
	Payload() CommandPayload
	Metadata() Metadata
	OccurredAt() time.Time
}

type command struct {
	Entity
	payload    CommandPayload
	metadata   Metadata
	occurredAt time.Time
}

var _ Command = (*command)(nil)

func NewCommand(name string, payload CommandPayload, options ...CommandOption) command {
	return newCommand(name, payload, options...)
}

func newCommand(name string, payload CommandPayload, options ...CommandOption) command {
	evt := command{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureCommand(&evt)
	}

	return evt
}

func (e command) CommandName() string     { return e.EntityName() }
func (e command) Payload() CommandPayload { return e.payload }
func (e command) Metadata() Metadata      { return e.metadata }
func (e command) OccurredAt() time.Time   { return e.occurredAt }
