package es

import (
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/ordering/internal/domain"
)

type orderRepository[T es.EventSourcedAggregate] struct {
	es.AggregateRepository[T]
}

func NewOrderRepository[T es.EventSourcedAggregate](
	registry registry.Registry, store es.AggregateStore,
) *orderRepository[T] {
	return &orderRepository[T]{
		AggregateRepository: es.NewAggregateRepository[T](
			domain.OrderAggregate,
			registry,
			store,
		),
	}
}

var _ domain.OrderRepository = (*orderRepository[*domain.Order])(nil)
