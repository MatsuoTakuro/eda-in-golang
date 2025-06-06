package ddd

type Metadata map[string]any

func (m Metadata) Set(key string, value any) {
	m[key] = value
}

func (m Metadata) Get(key string) any {
	return m[key]
}

func (m Metadata) Del(key string) {
	delete(m, key)
}

var (
	_ EventOption   = (*Metadata)(nil)
	_ CommandOption = (*Metadata)(nil)
	_ ReplyOption   = (*Metadata)(nil)
)

func (m Metadata) configureEvent(e *event) {
	for key, value := range m {
		e.metadata[key] = value
	}
}

func (m Metadata) configureCommand(c *command) {
	for key, value := range m {
		c.metadata[key] = value
	}
}

func (m Metadata) configureReply(r *reply) {
	for key, value := range m {
		r.metadata[key] = value
	}
}
