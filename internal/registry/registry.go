package registry

import (
	"sync"
)

// Registry provides a way to register types with their serializer, deserializer, and factory.
// It enables building, serializing, and deserializing registered types by key.
type Registry interface {
	// Serialize serializes v with a serializer registered by the given key
	Serialize(key string, v interface{}) ([]byte, error)
	MustSerialize(key string, v interface{}) []byte
	// Deserialize deserializes data with a deserializer registered by the given key
	Deserialize(key string, data []byte, options ...BuildOption) (interface{}, error)
	MustDeserialize(key string, data []byte, options ...BuildOption) interface{}
	// Build builds v from the given key that was registered with a factory function
	Build(key string, options ...BuildOption) (interface{}, error)
	MustBuild(key string, options ...BuildOption) interface{}
	register(key string, fn func() interface{}, s Serializer, d Deserializer, o []BuildOption) error
}

type Serializer func(v interface{}) ([]byte, error)
type Deserializer func(d []byte, v interface{}) error

type registry struct {
	registered map[string]registered
	mu         sync.RWMutex
}

type registered struct {
	factory      func() interface{}
	serializer   Serializer
	deserializer Deserializer
	options      []BuildOption
}

var _ Registry = (*registry)(nil)

// New creates a new registry.
// Note the registry is thread-safe and can be used concurrently as a singleton
func New() *registry {
	return &registry{
		registered: make(map[string]registered),
	}
}

func (r *registry) Serialize(key string, v interface{}) ([]byte, error) {
	reg, exists := r.registered[key]
	if !exists {
		return nil, unregisteredKey(key)
	}
	return reg.serializer(v)
}

func (r *registry) MustSerialize(key string, v interface{}) []byte {
	data, err := r.Serialize(key, v)
	if err != nil {
		panic(err)
	}
	return data
}

func (r *registry) Deserialize(key string, data []byte, options ...BuildOption) (interface{}, error) {
	v, err := r.Build(key, options...)
	if err != nil {
		return nil, err
	}

	err = r.registered[key].deserializer(data, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *registry) MustDeserialize(key string, data []byte, options ...BuildOption) interface{} {
	v, err := r.Deserialize(key, data, options...)
	if err != nil {
		panic(err)
	}
	return v
}

func (r *registry) Build(key string, options ...BuildOption) (interface{}, error) {
	reg, exists := r.registered[key]
	if !exists {
		return nil, unregisteredKey(key)
	}

	v := reg.factory()
	opts := append(r.registered[key].options, options...)

	for _, o := range opts {
		err := o(v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (r *registry) MustBuild(key string, options ...BuildOption) interface{} {
	v, err := r.Build(key, options...)
	if err != nil {
		panic(err)
	}
	return v
}

func (r *registry) register(key string, fn func() interface{}, s Serializer, d Deserializer, o []BuildOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registered[key]; exists {
		return alreadyRegisteredKey(key)
	}

	r.registered[key] = registered{
		factory:      fn,
		serializer:   s,
		deserializer: d,
		options:      o,
	}

	return nil
}
