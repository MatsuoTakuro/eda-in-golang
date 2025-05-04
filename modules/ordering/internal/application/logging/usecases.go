package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/modules/ordering/internal/application"
	"eda-in-golang/modules/ordering/internal/application/commands"
	"eda-in-golang/modules/ordering/internal/application/queries"
	"eda-in-golang/modules/ordering/internal/domain"
)

type Usecases struct {
	application.Usecases
	logger zerolog.Logger
}

var _ application.Usecases = (*Usecases)(nil)

func NewUsecases(usecases application.Usecases, logger zerolog.Logger) Usecases {
	return Usecases{
		Usecases: usecases,
		logger:   logger,
	}
}

func (a Usecases) CreateOrder(ctx context.Context, cmd commands.CreateOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.CreateOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.CreateOrder") }()
	return a.Usecases.CreateOrder(ctx, cmd)
}

func (a Usecases) CancelOrder(ctx context.Context, cmd commands.CancelOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.CancelOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.CancelOrder") }()
	return a.Usecases.CancelOrder(ctx, cmd)
}

func (a Usecases) ReadyOrder(ctx context.Context, cmd commands.ReadyOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.ReadyOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.ReadyOrder") }()
	return a.Usecases.ReadyOrder(ctx, cmd)
}

func (a Usecases) CompleteOrder(ctx context.Context, cmd commands.CompleteOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.CompleteOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.CompleteOrder") }()
	return a.Usecases.CompleteOrder(ctx, cmd)
}

func (a Usecases) GetOrder(ctx context.Context, query queries.GetOrder) (order *domain.Order, err error) {
	a.logger.Info().Msg("--> Ordering.GetOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.GetOrder") }()
	return a.Usecases.GetOrder(ctx, query)
}
