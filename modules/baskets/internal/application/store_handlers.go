package application

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/storespb"
)

type storeHandler[T ddd.Event] struct {
	logger zerolog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*storeHandler[ddd.Event])(nil)

func NewStoreHandler(logger zerolog.Logger) storeHandler[ddd.Event] {
	return storeHandler[ddd.Event]{
		logger: logger,
	}
}

func (h storeHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storespb.StoreParticipatingToggledEvent:
		return h.onStoreParticipationToggled(ctx, event)
	case storespb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}

	return nil
}

func (h storeHandler[T]) onStoreCreated(_ context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreCreated)
	h.logger.Debug().Msgf(`ID: %s, Name: "%s", Location: "%s"`, payload.GetId(), payload.GetName(), payload.GetLocation())
	return nil
}

func (h storeHandler[T]) onStoreParticipationToggled(_ context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreParticipationToggled)
	h.logger.Debug().Msgf(`ID: %s, Participating: %t`, payload.GetId(), payload.Participating)
	return nil
}

func (h storeHandler[T]) onStoreRebranded(_ context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreRebranded)
	h.logger.Debug().Msgf(`ID: %s, Name: "%s"`, payload.GetId(), payload.GetName())
	return nil
}
