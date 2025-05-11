package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
)

type eventHandler[T ddd.Event] struct {
	ddd.EventHandler[T]
	label  string
	logger zerolog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*eventHandler[ddd.Event])(nil)

func NewEventHandler[T ddd.Event](handlers ddd.EventHandler[T], label string, logger zerolog.Logger) eventHandler[T] {
	return eventHandler[T]{
		EventHandler: handlers,
		label:        label,
		logger:       logger,
	}
}

func (h eventHandler[T]) HandleEvent(ctx context.Context, event T) (err error) {
	h.logger.Info().Msgf("--> Ordering.%s.On(%s)", h.label, event.EventName())
	defer func() { h.logger.Info().Err(err).Msgf("<-- Ordering.%s.On(%s)", h.label, event.EventName()) }()
	return h.EventHandler.HandleEvent(ctx, event)
}
