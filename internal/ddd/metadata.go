package ddd

import "maps"

type Metadata map[string]any

var _ EventOption = (*Metadata)(nil)

func (m Metadata) Set(key string, value any) {
	m[key] = value
}

func (m Metadata) Get(key string) any {
	return m[key]
}

func (m Metadata) Del(key string) {
	delete(m, key)
}

func (m Metadata) apply(e *event) {
	maps.Copy(e.metadata, m)
}
