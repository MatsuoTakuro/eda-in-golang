package registry

import (
	"fmt"
)

type (
	unregisteredKey      string
	alreadyRegisteredKey string
)

func (key unregisteredKey) Error() string {
	return fmt.Sprintf("nothing has been registered with the key `%s`", string(key))
}

func (key alreadyRegisteredKey) Error() string {
	return fmt.Sprintf("something with the key `%s` has already been registered", string(key))
}
