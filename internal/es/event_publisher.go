package es

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type eventPublisher struct {
	AggregateStore
	publisher ddd.EventPublisher[ddd.AggregateEvent]
}

func WithEventPublisher(publisher ddd.EventPublisher[ddd.AggregateEvent]) AggregateStoreMiddleware {
	eventPublisher := eventPublisher{
		publisher: publisher,
	}

	return func(store AggregateStore) AggregateStore {
		eventPublisher.AggregateStore = store
		return eventPublisher
	}
}

func (p eventPublisher) Save(ctx context.Context, aggregate EventSourcedAggregate) error {
	if err := p.AggregateStore.Save(ctx, aggregate); err != nil {
		return err
	}
	return p.publisher.Publish(ctx, aggregate.Events()...)
}
