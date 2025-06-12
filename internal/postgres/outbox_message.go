package postgres

import "eda-in-golang/internal/am"

type outboxMessage struct {
	id      string
	name    string
	subject string
	data    []byte
}

var _ am.RawMessage = (*outboxMessage)(nil)

func (m outboxMessage) ID() string {
	return m.id
}

func (m outboxMessage) Subject() string {
	return m.subject
}

func (m outboxMessage) MessageName() string {
	return m.name
}

func (m outboxMessage) Data() []byte {
	return m.data
}
