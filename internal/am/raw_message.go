package am

type RawMessage interface {
	Message
	Data() []byte
}

type rawMessage struct {
	id      string
	name    string
	subject string
	data    []byte
}

var _ RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string          { return m.id }
func (m rawMessage) Subject() string     { return m.subject }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte        { return m.data }

type AckableRawMessage interface {
	AckableMessage
	Data() []byte
}
