package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/internal/domain"
	"eda-in-golang/modules/stores/storespb"
)

type storeHandler[T ddd.Event] struct {
	cache domain.StoreCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*storeHandler[ddd.Event])(nil)

func NewStoreHandler(cache domain.StoreCacheRepository) storeHandler[ddd.Event] {
	return storeHandler[ddd.Event]{
		cache: cache,
	}
}

func (h storeHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storespb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}

	return nil
}

func (h storeHandler[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreCreated)
	return h.cache.Add(ctx, payload.GetId(), payload.GetName(), payload.GetLocation())
}

func (h storeHandler[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreRebranded)
	return h.cache.Rename(ctx, payload.GetId(), payload.GetName())
}
