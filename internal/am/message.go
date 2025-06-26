package am

import (
	"time"

	"eda-in-golang/internal/ddd"
)

type Message interface {
	MessageBase
	Data() []byte
}

// MessageBase is a message interface that represents a message in the system.
// It can be an event, command, or any other type of message.
type MessageBase interface {
	// ID returns the unique identifier of the message.
	ddd.IDer
	// Subject returns the subject (topic) of the message.
	Subject() string
	// MessageName returns the type name of the message.
	MessageName() string
	Metadata() ddd.Metadata
	SentAt() time.Time
}

type message struct {
	id       string
	name     string
	subject  string
	data     []byte
	metadata ddd.Metadata
	sentAt   time.Time
}

var _ Message = (*message)(nil)

func (m message) ID() string             { return m.id }
func (m message) Subject() string        { return m.subject }
func (m message) MessageName() string    { return m.name }
func (m message) Data() []byte           { return m.data }
func (m message) Metadata() ddd.Metadata { return m.metadata }
func (m message) SentAt() time.Time      { return m.sentAt }

type IncomingMessage interface {
	IncomingMessageBase
	Data() []byte
}

// IncomingMessageBase represents a message that supports acknowledgment control.
type IncomingMessageBase interface {
	MessageBase
	ReceivedAt() time.Time
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}
