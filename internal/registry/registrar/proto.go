package registrar

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"

	"eda-in-golang/internal/registry"
)

type protoRegistrar struct {
	r registry.Registry
}

var _ registry.Registrar = (*protoRegistrar)(nil)
var protoT = reflect.TypeOf((*proto.Message)(nil)).Elem()

// NewProtoRegistrar creates a new registry registrar with protobuf serialization/deserialization.
func NewProtoRegistrar(r registry.Registry) *protoRegistrar {
	return &protoRegistrar{r: r}
}

func (c protoRegistrar) Register(v registry.Registrable, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message", v)
	}
	return registry.Register(c.r, v, c.serialize, c.deserialize, options)
}

func (c protoRegistrar) RegisterWithKey(key string, v interface{}, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message", v)
	}
	return registry.RegisterWithKey(c.r, key, v, c.serialize, c.deserialize, options)
}

func (protoRegistrar) serialize(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (protoRegistrar) deserialize(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
