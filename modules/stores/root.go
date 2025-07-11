package stores

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/es"
	evtstm "eda-in-golang/modules/stores/internal/es"

	"eda-in-golang/internal/jetstream"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
	"eda-in-golang/internal/system"

	"eda-in-golang/internal/tm"
	"eda-in-golang/modules/stores/internal/application"
	"eda-in-golang/modules/stores/internal/domain"
	"eda-in-golang/modules/stores/internal/grpc"
	"eda-in-golang/modules/stores/internal/handlers"
	"eda-in-golang/modules/stores/internal/logging"
	"eda-in-golang/modules/stores/internal/postgres"
	"eda-in-golang/modules/stores/internal/rest"
	"eda-in-golang/modules/stores/storespb"
)

func Root(ctx context.Context, svc system.Service) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := storespb.RegisterMessages(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	container.AddSingleton("logger", func(c di.Container) (any, error) {
		return svc.Logger(), nil
	})
	container.AddSingleton("stream", func(c di.Container) (any, error) {
		return jetstream.NewStream("stores", svc.Config().Nats.Stream, svc.JS(), c.Get("logger").(zerolog.Logger)), nil
	})
	container.AddSingleton("domainDispatcher", func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})
	container.AddSingleton("db", func(c di.Container) (any, error) {
		return svc.DB(), nil
	})
	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		return tm.NewOutboxProcessor(
			c.Get("stream").(am.RawMessageStream),
			pg.NewOutboxStore("stores.outbox", c.Get("db").(*sql.DB)),
		), nil
	})
	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		return db.Begin()
	})
	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		outboxStore := pg.NewOutboxStore("stores.outbox", tx)
		return am.WithRawMessageStreamMiddlewares(
			c.Get("stream").(am.RawMessageStream),
			tm.WithOutboxStream(outboxStore),
		), nil
	})
	container.AddScoped("eventStream", func(c di.Container) (any, error) {
		return am.NewEventStream(c.Get("registry").(registry.Registry), c.Get("txStream").(am.RawMessageStream)), nil
	})
	container.AddScoped("aggregateStore", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		reg := c.Get("registry").(registry.Registry)
		return es.AggregateStoreWithMiddleware(
			pg.NewEventStore("stores.events", tx, reg),
			pg.WithSnapshotStore("stores.snapshots", tx, reg),
		), nil
	})
	container.AddScoped("stores", func(c di.Container) (any, error) {
		return evtstm.NewStoreRepository[*domain.Store](
			c.Get("registry").(registry.Registry),
			c.Get("aggregateStore").(es.AggregateStore),
		), nil
	})
	container.AddScoped("products", func(c di.Container) (any, error) {
		return evtstm.NewProductRepository[*domain.Product](
			c.Get("registry").(registry.Registry),
			c.Get("aggregateStore").(es.AggregateStore),
		), nil
	})
	container.AddScoped("catalog", func(c di.Container) (any, error) {
		return postgres.NewCatalogRepository("stores.products", c.Get("tx").(*sql.Tx)), nil
	})
	container.AddScoped("mall", func(c di.Container) (any, error) {
		return postgres.NewMallRepository("stores.stores", c.Get("tx").(*sql.Tx)), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		return logging.LogApplicationAccess(
			application.New(
				c.Get("stores").(domain.StoreRepository),
				c.Get("products").(domain.ProductRepository),
				c.Get("catalog").(domain.CatalogRepository),
				c.Get("mall").(domain.MallRepository),
				c.Get("domainDispatcher").(ddd.EventDispatcher[ddd.Event]),
			),
			c.Get("logger").(zerolog.Logger),
		), nil
	})
	container.AddScoped("catalogHandlers", func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess(
			handlers.NewCatalogHandlers(c.Get("catalog").(domain.CatalogRepository)),
			"Catalog", c.Get("logger").(zerolog.Logger),
		), nil
	})
	container.AddScoped("mallHandlers", func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess(
			handlers.NewMallHandlers(c.Get("mall").(domain.MallRepository)),
			"Mall", c.Get("logger").(zerolog.Logger),
		), nil
	})
	container.AddScoped("domainEventHandlers", func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess(
			handlers.NewDomainEventHandlers(c.Get("eventStream").(am.EventStream)),
			"DomainEvents", c.Get("logger").(zerolog.Logger),
		), nil
	})

	// setup Driver adapters
	if err = grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.RegisterCatalogHandlersTx(container)
	handlers.RegisterMallHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)
	if err = storespb.RegisterAsyncAPI(svc.Mux()); err != nil {
		return err
	}
	startOutboxProcessor(ctx, container)

	return nil
}

func registrations(reg registry.Registry) (err error) {
	regtr := registrar.NewJsonRegistrar(reg)

	// Store
	if err = regtr.Register(domain.Store{}, func(v any) error {
		store := v.(*domain.Store)
		store.Aggregate = es.NewAggregate("", domain.StoreAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = regtr.Register(domain.StoreCreated{}); err != nil {
		return
	}
	if err = regtr.RegisterWithKey(domain.StoreParticipationEnabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return
	}
	if err = regtr.RegisterWithKey(domain.StoreParticipationDisabledEvent, domain.StoreParticipationToggled{}); err != nil {
		return
	}
	if err = regtr.Register(domain.StoreRebranded{}); err != nil {
		return
	}
	// store snapshots
	if err = regtr.RegisterWithKey(domain.StoreV1{}.SnapshotName(), domain.StoreV1{}); err != nil {
		return
	}

	// Product
	if err = regtr.Register(domain.Product{}, func(v any) error {
		store := v.(*domain.Product)
		store.Aggregate = es.NewAggregate("", domain.ProductAggregate)
		return nil
	}); err != nil {
		return
	}
	// product events
	if err = regtr.Register(domain.ProductAdded{}); err != nil {
		return
	}
	if err = regtr.Register(domain.ProductRebranded{}); err != nil {
		return
	}
	if err = regtr.RegisterWithKey(domain.ProductPriceIncreasedEvent, domain.ProductPriceChanged{}); err != nil {
		return
	}
	if err = regtr.RegisterWithKey(domain.ProductPriceDecreasedEvent, domain.ProductPriceChanged{}); err != nil {
		return
	}
	if err = regtr.Register(domain.ProductRemoved{}); err != nil {
		return
	}
	// product snapshots
	if err = regtr.RegisterWithKey(domain.ProductV1{}.SnapshotName(), domain.ProductV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, container di.Container) {
	outboxProcessor := container.Get("outboxProcessor").(tm.OutboxProcessor)
	logger := container.Get("logger").(zerolog.Logger)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("stores outbox processor encountered an error")
		}
	}()
}
