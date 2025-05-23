package baskets

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
	"eda-in-golang/modules/baskets/internal/application"
	"eda-in-golang/modules/baskets/internal/domain"
	"eda-in-golang/modules/baskets/internal/grpc"
	"eda-in-golang/modules/baskets/internal/handlers"
	"eda-in-golang/modules/baskets/internal/logging"
	"eda-in-golang/modules/baskets/internal/rest"
	"eda-in-golang/modules/stores/storespb"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono monolith.Server) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = registrations(reg); err != nil {
		return err
	}
	if err = storespb.Registrations(reg); err != nil {
		return err
	}
	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("baskets.events", mono.DB(), reg),
		es.WithEventPublisher(domainDispatcher),
		pg.WithSnapshotStore("baskets.snapshots", mono.DB(), reg),
	)
	baskets := es.NewAggregateRepository[*domain.Basket](domain.BasketAggregate, reg, aggregateStore)
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(baskets, stores, products, orders),
		mono.Logger(),
	)
	orderHandlers := logging.LogEventHandlerAccess(
		application.NewOrderHandlers(orders),
		"Order", mono.Logger(),
	)
	storeHandler := logging.LogEventHandlerAccess(
		application.NewStoreHandler(mono.Logger()),
		"Store", mono.Logger(),
	)
	productHandlers := logging.LogEventHandlerAccess(
		application.NewProductHandlers(mono.Logger()),
		"Product", mono.Logger(),
	)

	// setup Driver adapters
	if err = grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}
	handlers.RegisterOrderHandlers(orderHandlers, domainDispatcher)
	if err = handlers.RegisterStoreHandler(storeHandler, eventStream); err != nil {
		return err
	}
	if err = handlers.RegisterProductHandlers(productHandlers, eventStream); err != nil {
		return err
	}

	return
}

func registrations(reg registry.Registry) error {
	regtr := registrar.NewJsonRegistrar(reg)

	// Basket
	if err := regtr.Register(domain.Basket{}, func(v interface{}) error {
		basket := v.(*domain.Basket)
		basket.Items = make(map[string]domain.Item)
		return nil
	}); err != nil {
		return err
	}
	// basket events
	if err := regtr.Register(domain.BasketStarted{}); err != nil {
		return err
	}
	if err := regtr.Register(domain.BasketCanceled{}); err != nil {
		return err
	}
	if err := regtr.Register(domain.BasketCheckedOut{}); err != nil {
		return err
	}
	if err := regtr.Register(domain.BasketItemAdded{}); err != nil {
		return err
	}
	if err := regtr.Register(domain.BasketItemRemoved{}); err != nil {
		return err
	}
	// basket snapshots
	if err := regtr.RegisterWithKey(domain.BasketV1{}.SnapshotName(), domain.BasketV1{}); err != nil {
		return err
	}

	return nil
}
