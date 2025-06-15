package es

import (
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/stores/internal/domain"
)

type productRepository[T es.EventSourcedAggregate] struct {
	es.AggregateRepository[T]
}

func NewProductRepository[T es.EventSourcedAggregate](
	registry registry.Registry, store es.AggregateStore,
) *productRepository[T] {
	return &productRepository[T]{
		AggregateRepository: es.NewAggregateRepository[T](
			domain.ProductAggregate,
			registry,
			store,
		),
	}
}

var _ domain.ProductRepository = (*productRepository[*domain.Product])(nil)
