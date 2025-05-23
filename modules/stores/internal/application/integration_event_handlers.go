package application

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/internal/domain"
	"eda-in-golang/modules/stores/storespb"
)

type integrationEventHandler[T ddd.AggregateEvent] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*integrationEventHandler[ddd.AggregateEvent])(nil)

// NewIntegrationEventHandler creates a new integration event handler for the store aggregate.
// Unlike domain events that work only within the service, integration events are public, stable messages meant to be consumed by other modules or services.
func NewIntegrationEventHandler(publisher am.MessagePublisher[ddd.Event]) *integrationEventHandler[ddd.AggregateEvent] {
	return &integrationEventHandler[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func (h integrationEventHandler[T]) HandleEvent(ctx context.Context, event T) error {
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

func (h integrationEventHandler[T]) onStoreCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreCreated)
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreCreatedEvent, &storespb.StoreCreated{
			Id:       event.ID(),
			Name:     payload.Name,
			Location: payload.Location,
		}),
	)
}

func (h integrationEventHandler[T]) onStoreParticipationEnabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreParticipatingToggledEvent, &storespb.StoreParticipationToggled{
			Id:            event.ID(),
			Participating: true,
		}),
	)
}

func (h integrationEventHandler[T]) onStoreParticipationDisabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreParticipatingToggledEvent, &storespb.StoreParticipationToggled{
			Id:            event.ID(),
			Participating: false,
		}),
	)
}

func (h integrationEventHandler[T]) onStoreRebranded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreRebranded)
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreRebrandedEvent, &storespb.StoreRebranded{
			Id:   event.ID(),
			Name: payload.Name,
		}),
	)
}
