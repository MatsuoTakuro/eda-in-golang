package es

import (
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/stores/internal/domain"
)

type storeRepository[T es.EventSourcedAggregate] struct {
	es.AggregateRepository[T]
}

func NewStoreRepository[T es.EventSourcedAggregate](
	registry registry.Registry, store es.AggregateStore,
) *storeRepository[T] {
	return &storeRepository[T]{
		AggregateRepository: es.NewAggregateRepository[T](
			domain.StoreAggregate,
			registry,
			store,
		),
	}
}

var _ domain.StoreRepository = (*storeRepository[*domain.Store])(nil)
