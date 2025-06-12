package di

import (
	"context"
)

type containerContextKey struct{}

var containerKey = containerContextKey{}

// Get retrieves a dependency from the container that must be built beforehand.
func Get(ctx context.Context, key Key) any {
	cntr, ok := ctx.Value(containerKey).(*container)
	if !ok {
		panic("container does not exist on context")
	}

	return cntr.Get(key)
}
