package es

import (
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/baskets/internal/domain"
)

type bascketRepository[T es.EventSourcedAggregate] struct {
	es.AggregateRepository[T]
}

func NewBasketRepository[T es.EventSourcedAggregate](
	registry registry.Registry, store es.AggregateStore,
) *bascketRepository[T] {
	return &bascketRepository[T]{
		AggregateRepository: es.NewAggregateRepository[T](
			domain.BasketAggregate,
			registry,
			store,
		),
	}
}

var _ domain.BasketRepository = (*bascketRepository[*domain.Basket])(nil)
