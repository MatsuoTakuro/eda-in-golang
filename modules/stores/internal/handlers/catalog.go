package handlers

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/modules/stores/internal/domain"
)

type catalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

func NewCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ProductAddedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreasedEvent,
		domain.ProductPriceDecreasedEvent,
		domain.ProductRemovedEvent,
	)
}

func RegisterCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, "catalogHandlers").(ddd.EventHandler[ddd.Event])

		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(ddd.EventDispatcher[ddd.Event])

	RegisterCatalogHandlers(subscriber, handlers)
}

func (h catalogHandlers[T]) HandleEvent(ctx context.Context, event T) error {
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

func (h catalogHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.catalog.AddProduct(ctx, product.ID(), product.StoreID, product.Name, product.Description, product.SKU, product.Price)
}

func (h catalogHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.catalog.Rebrand(ctx, product.ID(), product.Name, product.Description)
}

func (h catalogHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	delta := event.Payload().(*domain.ProductPriceDelta)
	return h.catalog.UpdatePrice(ctx, delta.Product.ID(), delta.Delta)
}

func (h catalogHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	delta := event.Payload().(*domain.ProductPriceDelta)
	return h.catalog.UpdatePrice(ctx, delta.Product.ID(), delta.Delta)
}

func (h catalogHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.catalog.RemoveProduct(ctx, product.ID())
}
