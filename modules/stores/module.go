package stores

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/monolith"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
	"eda-in-golang/modules/stores/internal/application"
	"eda-in-golang/modules/stores/internal/domain"
	"eda-in-golang/modules/stores/internal/grpc"
	"eda-in-golang/modules/stores/internal/handlers"
	"eda-in-golang/modules/stores/internal/logging"
	"eda-in-golang/modules/stores/internal/postgres"
	"eda-in-golang/modules/stores/internal/rest"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, mono monolith.Server) error {
	// setup Driven adapters
	reg := registry.New()
	err := registrations(reg)
	if err != nil {
		return err
	}
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("stores.events", mono.DB(), reg),
		es.WithEventPublisher(domainDispatcher),
		pg.WithSnapshotStore("stores.snapshots", mono.DB(), reg),
	)
	stores := es.NewAggregateRepository[*domain.Store](domain.StoreAggregate, reg, aggregateStore)
	products := es.NewAggregateRepository[*domain.Product](domain.ProductAggregate, reg, aggregateStore)
	catalog := postgres.NewCatalogRepository("stores.products", mono.DB())
	mall := postgres.NewMallRepository("stores.stores", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(stores, products, catalog, mall),
		mono.Logger(),
	)
	catalogHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewCatalogHandlers(catalog),
		"Catalog", mono.Logger(),
	)
	mallHandlers := logging.LogEventHandlerAccess[ddd.AggregateEvent](
		application.NewMallHandlers(mall),
		"Mall", mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}
	handlers.RegisterCatalogHandlers(catalogHandlers, domainDispatcher)
	handlers.RegisterMallHandlers(mallHandlers, domainDispatcher)

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := registrar.NewJsonRegistrar(reg)

	// Store
	if err = serde.Register(domain.Store{}, func(v any) error {
		store := v.(*domain.Store)
		store.Aggregate = es.NewAggregate("", domain.StoreAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = serde.Register(domain.StoreCreated{}); err != nil {
		return
	}
	if err = serde.RegisterWithKey(domain.StoreParticipationEnabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return
	}
	if err = serde.RegisterWithKey(domain.StoreParticipationDisabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return
	}
	if err = serde.Register(domain.StoreRebranded{}); err != nil {
		return
	}
	// store snapshots
	if err = serde.Register(domain.StoreV1{}); err != nil {
		return
	}

	// Product
	if err = serde.Register(domain.Product{}, func(v any) error {
		store := v.(*domain.Product)
		store.Aggregate = es.NewAggregate("", domain.ProductAggregate)
		return nil
	}); err != nil {
		return
	}
	// product events
	if err = serde.Register(domain.ProductAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.ProductRebranded{}); err != nil {
		return
	}
	if err = serde.RegisterWithKey(domain.ProductPriceIncreasedEvent, domain.ProductPriceChanged{}); err != nil {
		return
	}
	if err = serde.RegisterWithKey(domain.ProductPriceDecreasedEvent, domain.ProductPriceChanged{}); err != nil {
		return
	}
	if err = serde.Register(domain.ProductRemoved{}); err != nil {
		return
	}
	// product snapshots
	if err = serde.Register(domain.ProductV1{}); err != nil {
		return
	}

	return
}
