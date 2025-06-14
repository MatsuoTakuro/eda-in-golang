package application

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/depot/internal/domain"
	"eda-in-golang/modules/stores/storespb"
)

type integrationHandler[T ddd.Event] struct {
	stores   domain.StoreCacheRepository
	products domain.ProductCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandler[ddd.Event])(nil)

// NewIntegrationEventHandler creates a handler that handle integration events coming from other modules.
func NewIntegrationEventHandler(
	stores domain.StoreCacheRepository,
	products domain.ProductCacheRepository,
) integrationHandler[ddd.Event] {
	return integrationHandler[ddd.Event]{
		stores:   stores,
		products: products,
	}
}

func (h integrationHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storespb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	case storespb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case storespb.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case storespb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}

	return nil
}

func (h integrationHandler[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreCreated)
	return h.stores.Add(ctx, payload.GetId(), payload.GetName(), payload.GetLocation())
}

func (h integrationHandler[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreRebranded)
	return h.stores.Rename(ctx, payload.GetId(), payload.GetName())
}

func (h integrationHandler[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductAdded)
	return h.products.Add(ctx, payload.GetId(), payload.GetStoreId(), payload.GetName())
}

func (h integrationHandler[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRebranded)
	return h.products.Rebrand(ctx, payload.GetId(), payload.GetName())
}

func (h integrationHandler[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRemoved)
	return h.products.Remove(ctx, payload.GetId())
}
