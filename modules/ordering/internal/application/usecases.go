package application

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/ordering/internal/application/commands"
	"eda-in-golang/modules/ordering/internal/application/queries"
	"eda-in-golang/modules/ordering/internal/domain/infra"
)

type Usecases interface {
	commands.Commands
	queries.Queries
}

type usecases struct {
	commands.Commands
	queries.Queries
}

var _ Usecases = (*usecases)(nil)

func NewUsecases(
	orderRepo infra.OrderRepository,
	shoppingClient infra.ShoppingClient,
	publisher ddd.EventPublisher[ddd.Event],
) Usecases {
	return &usecases{
		Commands: commands.New(orderRepo, shoppingClient, publisher),
		Queries:  queries.New(orderRepo),
	}
}
