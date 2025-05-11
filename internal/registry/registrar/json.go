package registrar

import (
	"encoding/json"

	"eda-in-golang/internal/registry"
)

type jsonRegistrar struct {
	r registry.Registry
}

var _ registry.Registrar = (*jsonRegistrar)(nil)

// NewJsonRegistrar creates a new registry registrar with JSON serialization/deserialization.
func NewJsonRegistrar(r registry.Registry) *jsonRegistrar {
	return &jsonRegistrar{r: r}
}

func (c jsonRegistrar) Register(v registry.Registrable, options ...registry.BuildOption) error {
	return registry.Register(c.r, v, c.serialize, c.deserialize, options)
}

func (c jsonRegistrar) RegisterWithKey(key string, v interface{}, options ...registry.BuildOption) error {
	return registry.RegisterWithKey(c.r, key, v, c.serialize, c.deserialize, options)
}

func (jsonRegistrar) serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonRegistrar) deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
