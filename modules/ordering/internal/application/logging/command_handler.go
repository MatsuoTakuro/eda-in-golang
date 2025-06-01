package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
)

type commandHandler[T ddd.Command] struct {
	ddd.CommandHandler[T]
	label  string
	logger zerolog.Logger
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandler[ddd.Command])(nil)

func NewCommandHandler[T ddd.Command](
	handler ddd.CommandHandler[T],
	label string,
	logger zerolog.Logger,
) commandHandler[T] {
	return commandHandler[T]{
		CommandHandler: handler,
		label:          label,
		logger:         logger,
	}
}

func (h commandHandler[T]) HandleCommand(ctx context.Context, command T) (reply ddd.Reply, err error) {
	h.logger.Info().Msgf("--> Ordering.%s.On(%s)", h.label, command.CommandName())
	defer func() { h.logger.Info().Err(err).Msgf("<-- Ordering.%s.On(%s)", h.label, command.CommandName()) }()
	return h.CommandHandler.HandleCommand(ctx, command)
}
