package depot

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/depot/internal/application"
	"eda-in-golang/modules/depot/internal/grpc"
	"eda-in-golang/modules/depot/internal/handlers"
	"eda-in-golang/modules/depot/internal/logging"
	"eda-in-golang/modules/depot/internal/postgres"
	"eda-in-golang/modules/depot/internal/rest"
	"eda-in-golang/modules/stores/storespb"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Server) error {
	// setup Driven adapters
	reg := registry.New()
	if err := storespb.Registrations(reg); err != nil {
		return err
	}
	eventStream := am.NewEventStream(reg, jetstream.NewStream("depot", mono.Config().Nats.Stream, mono.JS()))
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	shoppingLists := postgres.NewShoppingListRepository("depot.shopping_lists", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	stores := postgres.NewStoreCacheRepository("depot.stores_cache", mono.DB(), grpc.NewStoreRepository(conn))
	products := postgres.NewProductCacheRepository("depot.products_cache", mono.DB(), grpc.NewProductRepository(conn))
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(shoppingLists, stores, products, domainDispatcher),
		mono.Logger(),
	)
	orderHandler := logging.LogEventHandlerAccess(
		application.NewOrderHandler(orders),
		"Order", mono.Logger(),
	)
	storeHandler := logging.LogEventHandlerAccess(
		application.NewStoreHandler(stores),
		"Store", mono.Logger(),
	)
	productHandler := logging.LogEventHandlerAccess(
		application.NewProductHandler(products),
		"Product", mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.Register(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}
	handlers.SubscribeDomainEventsForOrder(orderHandler, domainDispatcher)
	if err = handlers.SubscribeStoreIntegrationEvents(storeHandler, eventStream); err != nil {
		return err
	}
	if err = handlers.SubscribeProductIntegrationEvents(productHandler, eventStream); err != nil {
		return err
	}

	return nil
}
