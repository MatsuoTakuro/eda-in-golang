package handlers

import (
	"context"
	"encoding/json"
	"testing"

	v4 "github.com/pact-foundation/pact-go/v2/message/v4"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
	"eda-in-golang/modules/depot/internal/application"
	"eda-in-golang/modules/depot/internal/domain"
	"eda-in-golang/modules/stores/storespb"
)

func TestStoresConsumer(t *testing.T) {
	type mocks struct {
		stores   *domain.MockStoreCacheRepository
		products *domain.MockProductCacheRepository
	}

	type rawEvent struct {
		Name    string         `json:"name"`
		Payload map[string]any `json:"payload"`
	}

	reg := registry.New()
	err := storespb.RegisterMessagesWithRegistrar(registrar.NewJsonRegistrar(reg))
	if err != nil {
		t.Fatal(err)
	}

	pact, err := v4.NewAsynchronousPact(v4.Config{
		Provider: "stores-pub",
		Consumer: "depot-sub",
		PactDir:  "./pacts",
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		given   []models.ProviderState
		content rawEvent
		on      func(m mocks)
	}{
		"a StoreCreated message": {
			content: rawEvent{
				Name: storespb.StoreCreatedEvent,
				Payload: map[string]any{
					"id":       "store-id",
					"name":     "NewStore",
					"location": "NewLocation",
				},
			},
			on: func(m mocks) {
				m.stores.On("Add", mock.Anything, "store-id", "NewStore", "NewLocation").Return(nil)
			},
		},
		"a StoreRebranded message": {
			content: rawEvent{
				Name: storespb.StoreRebrandedEvent,
				Payload: map[string]any{
					"id":   "store-id",
					"name": "RebrandedStore",
				},
			},
			on: func(m mocks) {
				m.stores.On("Rename", mock.Anything, "store-id", "RebrandedStore").Return(nil)
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := mocks{
				stores:   domain.NewMockStoreCacheRepository(t),
				products: domain.NewMockProductCacheRepository(t),
			}
			if tc.on != nil {
				tc.on(m)
			}
			handlers := application.NewIntegrationEventHandler(m.stores, m.products)
			msgHandlerFn := func(contents v4.MessageContents) error {
				event := contents.Content.(*rawEvent)

				data, err := json.Marshal(event.Payload)
				if err != nil {
					return err
				}
				payload := reg.MustDeserialize(event.Name, data)

				return handlers.HandleEvent(context.Background(), ddd.NewEvent(event.Name, payload))
			}

			message := pact.AddAsynchronousMessage()
			for _, given := range tc.given {
				message = message.GivenWithParameter(given)
			}

			assert.NoError(t, message.
				ExpectsToReceive(name).
				WithJSONContent(tc.content).
				AsType(&rawEvent{}).
				ConsumedBy(msgHandlerFn).
				Verify(t))
		})
	}
}
