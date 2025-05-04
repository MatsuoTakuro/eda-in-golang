package eventhandlers

import (
	"context"

	"eda-in-golang/internal/ddd"
)

/*
OrderXXXEventHandlers define reactions to domain events in the application layer.
Each method must match the `ddd.EventHandler` signature to allow direct subscription and trigger side effects like notifications.
*/

type OnOrderCreatedEventHandler interface {
	OnOrderCreated(ctx context.Context, event ddd.Event) error
}

type OnOrderReadiedEventHandler interface {
	OnOrderReadied(ctx context.Context, event ddd.Event) error
}

type OnOrderCanceledEventHandler interface {
	OnOrderCanceled(ctx context.Context, event ddd.Event) error
}

type OnOrderCompletedEventHandler interface {
	OnOrderCompleted(ctx context.Context, event ddd.Event) error
}
