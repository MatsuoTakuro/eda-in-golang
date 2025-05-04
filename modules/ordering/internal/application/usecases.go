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

func NewUsecases(orderRepo infra.OrderRepository,
	customerClient infra.CustomerClient, paymentClient infra.PaymentClient, shoppingClient infra.ShoppingClient,
	domainPublisher ddd.EventPublisher,
) Usecases {
	return &usecases{
		Commands: commands.New(orderRepo, customerClient, paymentClient, shoppingClient, domainPublisher),
		Queries:  queries.New(orderRepo),
	}
}
