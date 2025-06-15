package ddd

import (
	"time"

	"github.com/google/uuid"
)

type ReplyOption interface {
	configureReply(*reply)
}

type ReplyPayload any

type Reply interface {
	ID() string
	ReplyName() string
	Payload() ReplyPayload
	Metadata() Metadata
	OccurredAt() time.Time
}

type reply struct {
	Entity
	payload    ReplyPayload
	metadata   Metadata
	occurredAt time.Time
}

var _ Reply = (*reply)(nil)

func NewReply(name string, payload ReplyPayload, options ...ReplyOption) reply {
	return newReply(name, payload, options...)
}

func newReply(name string, payload ReplyPayload, options ...ReplyOption) reply {
	rep := reply{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureReply(&rep)
	}

	return rep
}

func (e reply) ReplyName() string     { return e.EntityName() }
func (e reply) Payload() ReplyPayload { return e.payload }
func (e reply) Metadata() Metadata    { return e.metadata }
func (e reply) OccurredAt() time.Time { return e.occurredAt }
