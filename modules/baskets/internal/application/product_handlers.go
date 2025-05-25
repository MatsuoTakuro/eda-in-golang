package application

import (
	"context"

	"eda-in-golang/modules/baskets/internal/domain"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/storespb"
)

type productHandler[T ddd.Event] struct {
	cache domain.ProductCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*productHandler[ddd.Event])(nil)

func NewProductHandler(cache domain.ProductCacheRepository) productHandler[ddd.Event] {
	return productHandler[ddd.Event]{
		cache: cache,
	}
}

func (h productHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case storespb.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case storespb.ProductPriceIncreasedEvent, storespb.ProductPriceDecreasedEvent:
		return h.onProductPriceChanged(ctx, event)
	case storespb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}

	return nil
}

func (h productHandler[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductAdded)
	return h.cache.Add(ctx, payload.GetId(), payload.GetStoreId(), payload.GetName(), payload.GetPrice())
}

func (h productHandler[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRebranded)
	return h.cache.Rebrand(ctx, payload.GetId(), payload.GetName())
}

func (h productHandler[T]) onProductPriceChanged(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductPriceChanged)
	return h.cache.UpdatePrice(ctx, payload.GetId(), payload.GetDelta())
}

func (h productHandler[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRemoved)
	return h.cache.Remove(ctx, payload.GetId())
}
