package am

import (
	"time"

	"eda-in-golang/internal/ddd"
)

type CommandMessage interface {
	AckableMessage
	ddd.Command
}

type commandMessage struct {
	id         string
	name       string
	payload    ddd.CommandPayload
	metadata   ddd.Metadata
	occurredAt time.Time
	msg        AckableMessage
}

var _ CommandMessage = (*commandMessage)(nil)

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Payload() ddd.CommandPayload { return c.payload }
func (c commandMessage) Metadata() ddd.Metadata      { return c.metadata }
func (c commandMessage) OccurredAt() time.Time       { return c.occurredAt }
func (c commandMessage) Subject() string             { return c.msg.Subject() }
func (c commandMessage) MessageName() string         { return c.msg.MessageName() }
func (c commandMessage) Ack() error                  { return c.msg.Ack() }
func (c commandMessage) NAck() error                 { return c.msg.NAck() }
func (c commandMessage) Extend() error               { return c.msg.Extend() }
func (c commandMessage) Kill() error                 { return c.msg.Kill() }
