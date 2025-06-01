package search

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/registry"
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

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Server) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = orderingpb.RegisterMessages(reg); err != nil {
		return err
	}
	if err = customerspb.RegisterMessages(reg); err != nil {
		return err
	}
	if err = storespb.RegisterMessages(reg); err != nil {
		return err
	}
	eventStream := am.NewEventStream(reg, jetstream.NewStream("search", mono.Config().Nats.Stream, mono.JS()))
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	customers := postgres.NewCustomerCacheRepository("search.customers_cache", mono.DB(), grpc.NewCustomerRepository(conn))
	stores := postgres.NewStoreCacheRepository("search.stores_cache", mono.DB(), grpc.NewStoreRepository(conn))
	products := postgres.NewProductCacheRepository("search.products_cache", mono.DB(), grpc.NewProductRepository(conn))
	orders := postgres.NewOrderRepository("search.orders", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(orders),
		mono.Logger(),
	)
	orderHandlers := logging.LogEventHandlerAccess(
		application.NewOrderHandlers(orders, customers, stores, products),
		"Order", mono.Logger(),
	)
	customerHandlers := logging.LogEventHandlerAccess(
		application.NewCustomerHandlers(customers),
		"Customer", mono.Logger(),
	)
	storeHandlers := logging.LogEventHandlerAccess(
		application.NewStoreHandlers(stores),
		"Store", mono.Logger(),
	)
	productHandlers := logging.LogEventHandlerAccess(
		application.NewProductHandlers(products),
		"Product", mono.Logger(),
	)

	// setup Driver adapters
	if err = grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	if err = handlers.SubscribeOrderIntegrationEvents(
		orderHandlers, eventStream,
	); err != nil {
		return err
	}

	if err = handlers.SubscribeCustomerIntegrationEvents(
		customerHandlers, eventStream,
	); err != nil {
		return err
	}

	if err = handlers.SubscribeStoreIntegrationEvents(
		storeHandlers, eventStream,
	); err != nil {
		return err
	}

	if err = handlers.SubscribeProductIntegrationEvents(
		productHandlers, eventStream,
	); err != nil {
		return err
	}

	return nil
}
