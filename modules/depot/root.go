package depot

import (
	"context"
	"database/sql"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/jetstream"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/system"
	"eda-in-golang/internal/tm"
	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/depot/internal/application"
	dep "eda-in-golang/modules/depot/internal/di"
	"eda-in-golang/modules/depot/internal/domain"
	"eda-in-golang/modules/depot/internal/grpc"
	"eda-in-golang/modules/depot/internal/handlers"
	"eda-in-golang/modules/depot/internal/logging"
	"eda-in-golang/modules/depot/internal/postgres"
	"eda-in-golang/modules/depot/internal/rest"
	"eda-in-golang/modules/stores/storespb"

	"github.com/rs/zerolog"
)

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()

	// setup Driven adapters
	container.AddSingleton(di.Registry, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := storespb.RegisterMessages(reg); err != nil {
			return nil, err
		}
		if err := depotpb.RegisterMessages(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	container.AddSingleton(di.Logger, func(c di.Container) (any, error) {
		return svc.Logger(), nil
	})
	container.AddSingleton(di.Stream, func(c di.Container) (any, error) {
		return jetstream.NewStream(
				"depot", svc.Config().Nats.Stream, svc.JS(), c.Get(di.Logger).(zerolog.Logger),
			),
			nil
	})
	container.AddSingleton(di.DomainDispatcher, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})
	container.AddSingleton(di.DB, func(c di.Container) (any, error) {
		return svc.DB(), nil
	})
	container.AddSingleton("storesConn", func(c di.Container) (any, error) {
		return grpc.Dial(ctx, svc.Config().Rpc.Service("STORES"))
	})
	container.AddSingleton(di.OutboxProcessor, func(c di.Container) (any, error) {
		return tm.NewOutboxProcessor(
			c.Get(di.Stream).(am.RawMessageStream),
			pg.NewOutboxStore("depot.outbox", c.Get(di.DB).(*sql.DB)),
		), nil
	})
	container.AddScoped(di.TX, func(c di.Container) (any, error) {
		db := c.Get(di.DB).(*sql.DB)
		return db.Begin()
	})
	container.AddScoped(di.TXStream, func(c di.Container) (any, error) {
		tx := c.Get(di.TX).(*sql.Tx)
		outboxStore := pg.NewOutboxStore("depot.outbox", tx)
		return am.WithRawMessageStreamMiddlewares(
			c.Get(di.Stream).(am.RawMessageStream),
			tm.WithOutboxStream(outboxStore),
		), nil
	})
	container.AddScoped(di.EventStream, func(c di.Container) (any, error) {
		return am.NewEventStream(
			c.Get(di.Registry).(registry.Registry),
			c.Get(di.TXStream).(am.RawMessageStream),
		), nil
	})
	container.AddScoped(di.CommandStream, func(c di.Container) (any, error) {
		return am.NewCommandStream(
			c.Get(di.Registry).(registry.Registry),
			c.Get(di.TXStream).(am.RawMessageStream),
		), nil
	})
	container.AddScoped(di.ReplyStream, func(c di.Container) (any, error) {
		return am.NewReplyStream(
			c.Get(di.Registry).(registry.Registry), c.Get(di.TXStream).(am.RawMessageStream),
		), nil
	})
	container.AddScoped(di.InboxMiddleware, func(c di.Container) (any, error) {
		tx := c.Get(di.TX).(*sql.Tx)
		inboxStore := pg.NewInboxStore("depot.inbox", tx)
		return tm.WithInboxHandler(inboxStore), nil
	})
	container.AddScoped(dep.ShoppingLists, func(c di.Container) (any, error) {
		return postgres.NewShoppingListRepository(
			"depot.shopping_lists",
			c.Get(di.TX).(*sql.Tx),
		), nil
	})
	container.AddScoped(dep.Stores, func(c di.Container) (any, error) {
		return postgres.NewStoreCacheRepository(
			"depot.stores_cache",
			c.Get(di.TX).(*sql.Tx),
			grpc.NewStoreRepository(c.Get("storesConn").(*grpc.ClientConn)),
		), nil
	})
	container.AddScoped(dep.Products, func(c di.Container) (any, error) {
		return postgres.NewProductCacheRepository(
			"depot.products_cache",
			c.Get(di.TX).(*sql.Tx),
			grpc.NewProductRepository(c.Get("storesConn").(*grpc.ClientConn)),
		), nil
	})

	// setup application
	container.AddScoped(di.Application, func(c di.Container) (any, error) {
		return logging.LogApplicationAccess(
			application.New(
				c.Get(dep.ShoppingLists).(domain.ShoppingListRepository),
				c.Get(dep.Stores).(domain.StoreCacheRepository),
				c.Get(dep.Products).(domain.ProductCacheRepository),
				c.Get(di.DomainDispatcher).(ddd.EventDispatcher[ddd.AggregateEvent]),
			),
			c.Get(di.Logger).(zerolog.Logger),
		), nil
	})
	container.AddScoped(di.DomainEventHandler, func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess(
			application.NewDomainEventHandler(c.Get(di.EventStream).(am.EventStream)),
			"DomainEvents", c.Get(di.Logger).(zerolog.Logger),
		), nil
	})
	container.AddScoped(di.IntegrationEventHandler, func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess(
			application.NewIntegrationEventHandler(
				c.Get(dep.Stores).(domain.StoreCacheRepository),
				c.Get(dep.Products).(domain.ProductCacheRepository),
			),
			"IntegrationEvents", c.Get(di.Logger).(zerolog.Logger),
		), nil
	})
	container.AddScoped(di.CommandHandler, func(c di.Container) (any, error) {
		return logging.LogCommandHandlerAccess[ddd.Command](
			application.NewCommandHandler(c.Get(di.Application).(application.App)),
			"Commands", c.Get(di.Logger).(zerolog.Logger),
		), nil
	})

	// setup Driver adapters
	if err := grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.SubscribeDomainEvents(container)
	if err := handlers.SubscribeIntegrationEvents(container); err != nil {
		return err
	}
	if err := handlers.SubscribeCommands(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, container)

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	outboxProcessor := container.Get(di.OutboxProcessor).(tm.OutboxProcessor)
	logger := container.Get(di.Logger).(zerolog.Logger)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("depot outbox processor encountered an error")
		}
	}()
}
