package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/modules/stores/internal/domain"
)

func RegisterIntegrationEventHandler(eventHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(eventHandlers,
		domain.StoreCreatedEvent,
		domain.StoreParticipationEnabledEvent,
		domain.StoreParticipationDisabledEvent,
		domain.StoreRebrandedEvent,
	)
}
