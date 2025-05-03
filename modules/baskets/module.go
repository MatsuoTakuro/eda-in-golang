package baskets

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/modules/baskets/internal/application"
	"eda-in-golang/modules/baskets/internal/grpc"
	"eda-in-golang/modules/baskets/internal/handlers"
	"eda-in-golang/modules/baskets/internal/logging"
	"eda-in-golang/modules/baskets/internal/postgres"
	"eda-in-golang/modules/baskets/internal/rest"
)

type Module struct{}

var _ monolith.Module = (*Module)(nil)

func (m *Module) Startup(ctx context.Context, srv monolith.Server) (err error) {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	baskets := postgres.NewBasketRepository("baskets.baskets", srv.DB())
	conn, err := grpc.Dial(ctx, srv.Config().Rpc.Address())
	if err != nil {
		return err
	}
	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(baskets, stores, products, orders, domainDispatcher),
		srv.Logger(),
	)
	orderHandlers := logging.LogDomainEventHandlerAccess(
		application.NewOrderHandlers(orders),
		srv.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(app, srv.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, srv.Mux(), srv.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(srv.Mux()); err != nil {
		return err
	}
	handlers.RegisterOrderHandlers(orderHandlers, domainDispatcher)

	return
}
