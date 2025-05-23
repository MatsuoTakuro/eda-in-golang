package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/internal/domain"
)

type catalogHandler[T ddd.AggregateEvent] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*catalogHandler[ddd.AggregateEvent])(nil)

func NewCatalogHandler(catalog domain.CatalogRepository) *catalogHandler[ddd.AggregateEvent] {
	return &catalogHandler[ddd.AggregateEvent]{
		catalog: catalog,
	}
}

func (h catalogHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case domain.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case domain.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)
	case domain.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)
	case domain.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}
	return nil
}

func (h catalogHandler[T]) onProductAdded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductAdded)
	return h.catalog.AddProduct(ctx, event.AggregateID(), payload.StoreID, payload.Name, payload.Description, payload.SKU, payload.Price)
}

func (h catalogHandler[T]) onProductRebranded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductRebranded)
	return h.catalog.Rebrand(ctx, event.AggregateID(), payload.Name, payload.Description)
}

func (h catalogHandler[T]) onProductPriceIncreased(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductPriceChanged)
	return h.catalog.UpdatePrice(ctx, event.AggregateID(), payload.Delta)
}

func (h catalogHandler[T]) onProductPriceDecreased(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductPriceChanged)
	return h.catalog.UpdatePrice(ctx, event.AggregateID(), payload.Delta)
}

func (h catalogHandler[T]) onProductRemoved(ctx context.Context, event ddd.AggregateEvent) error {
	return h.catalog.RemoveProduct(ctx, event.AggregateID())
}
