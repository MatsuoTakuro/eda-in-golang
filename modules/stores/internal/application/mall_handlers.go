package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/internal/domain"
)

type mallHandler[T ddd.AggregateEvent] struct {
	mall domain.MallRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*mallHandler[ddd.AggregateEvent])(nil)

func NewMallHandler(mall domain.MallRepository) *mallHandler[ddd.AggregateEvent] {
	return &mallHandler[ddd.AggregateEvent]{
		mall: mall,
	}
}

func (h mallHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case domain.StoreParticipationEnabledEvent:
		return h.onStoreParticipationEnabled(ctx, event)
	case domain.StoreParticipationDisabledEvent:
		return h.onStoreParticipationDisabled(ctx, event)
	case domain.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}
	return nil
}

func (h mallHandler[T]) onStoreCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreCreated)
	return h.mall.AddStore(ctx, event.AggregateID(), payload.Name, payload.Location)
}

func (h mallHandler[T]) onStoreParticipationEnabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.mall.SetStoreParticipation(ctx, event.AggregateID(), true)
}

func (h mallHandler[T]) onStoreParticipationDisabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.mall.SetStoreParticipation(ctx, event.AggregateID(), false)
}

func (h mallHandler[T]) onStoreRebranded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreRebranded)
	return h.mall.RenameStore(ctx, event.AggregateID(), payload.Name)
}
