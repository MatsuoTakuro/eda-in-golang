package registry

import "reflect"

// Registrar registers types with their serialization and deserialization logic.
type Registrar interface {
	Register(v Registrable, options ...BuildOption) error
	RegisterWithKey(key string, v interface{}, options ...BuildOption) error
}

type Registrable interface{ Key() string }

// Register registers v to a registry
func Register(reg Registry, v Registrable, s Serializer, d Deserializer, os []BuildOption) error {
	var key string

	t := reflect.TypeOf(v)

	switch {
	// accept: (*T)(nil)
	case t.Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil():
		key = reflect.New(t).Interface().(Registrable).Key()
	// accept: *T{}, T{}
	default:
		key = v.Key()
	}

	return RegisterWithKey(reg, key, v, s, d, os)
}

// RegisterWithKey registers v to a registry with a specific key
func RegisterWithKey(reg Registry, key string, v interface{}, s Serializer, d Deserializer, os []BuildOption) error {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return reg.register(key, func() interface{} {
		return reflect.New(t).Interface()
	}, s, d, os)
}
