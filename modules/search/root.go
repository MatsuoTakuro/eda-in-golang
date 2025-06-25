package search

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/jetstream"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/system"
	"eda-in-golang/internal/tm"
	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/ordering/orderingpb"
	"eda-in-golang/modules/search/internal/application"
	"eda-in-golang/modules/search/internal/grpc"
	"eda-in-golang/modules/search/internal/handlers"
	"eda-in-golang/modules/search/internal/logging"
	"eda-in-golang/modules/search/internal/postgres"
	"eda-in-golang/modules/search/internal/rest"
	"eda-in-golang/modules/stores/storespb"
)

func Root(ctx context.Context, svc system.Service) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := customerspb.Registrations(reg); err != nil {
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
		return jetstream.NewStream("search", svc.Config().Nats.Stream, svc.JS(), c.Get("logger").(zerolog.Logger)), nil
	})
	container.AddSingleton("db", func(c di.Container) (any, error) {
		return svc.DB(), nil
	})
	container.AddSingleton("storesConn", func(c di.Container) (any, error) {
		return grpc.Dial(ctx, svc.Config().Rpc.Service("STORES"))
	})
	container.AddSingleton("customersConn", func(c di.Container) (any, error) {
		return grpc.Dial(ctx, svc.Config().Rpc.Service("CUSTOMERS"))
	})
	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		return db.Begin()
	})
	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		inboxStore := pg.NewInboxStore("search.inbox", tx)
		return tm.WithInboxHandler(inboxStore), nil
	})
	container.AddScoped("customers", func(c di.Container) (any, error) {
		return postgres.NewCustomerCacheRepository(
			"search.customers_cache",
			c.Get("tx").(*sql.Tx),
			grpc.NewCustomerRepository(c.Get("customersConn").(*grpc.ClientConn)),
		), nil
	})
	container.AddScoped("stores", func(c di.Container) (any, error) {
		return postgres.NewStoreCacheRepository(
			"search.stores_cache",
			c.Get("tx").(*sql.Tx),
			grpc.NewStoreRepository(c.Get("storesConn").(*grpc.ClientConn)),
		), nil
	})
	container.AddScoped("products", func(c di.Container) (any, error) {
		return postgres.NewProductCacheRepository(
			"search.products_cache",
			c.Get("tx").(*sql.Tx),
			grpc.NewProductRepository(c.Get("storesConn").(*grpc.ClientConn)),
		), nil
	})
	container.AddScoped("orders", func(c di.Container) (any, error) {
		return postgres.NewOrderRepository("search.orders", c.Get("tx").(*sql.Tx)), nil
	})

	// setup application
	container.AddScoped("app", func(c di.Container) (any, error) {
		return logging.LogApplicationAccess(
			application.New(
				c.Get("orders").(application.OrderRepository),
			),
			c.Get("logger").(zerolog.Logger),
		), nil
	})
	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess[ddd.Event](
			handlers.NewIntegrationEventHandlers(
				c.Get("orders").(application.OrderRepository),
				c.Get("customers").(application.CustomerCacheRepository),
				c.Get("stores").(application.StoreCacheRepository),
				c.Get("products").(application.ProductCacheRepository),
			),
			"IntegrationEvents", c.Get("logger").(zerolog.Logger),
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
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}

	return nil
}
