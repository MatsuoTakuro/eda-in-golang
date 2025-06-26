package postgres

import (
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"time"
)

type outboxMessage struct {
	id       string
	name     string
	subject  string
	data     []byte
	metadata ddd.Metadata
	sentAt   time.Time
}

var _ am.Message = (*outboxMessage)(nil)

func (m outboxMessage) ID() string             { return m.id }
func (m outboxMessage) Subject() string        { return m.subject }
func (m outboxMessage) MessageName() string    { return m.name }
func (m outboxMessage) Data() []byte           { return m.data }
func (m outboxMessage) Metadata() ddd.Metadata { return m.metadata }
func (m outboxMessage) SentAt() time.Time      { return m.sentAt }
