package depot

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/modules/depot/internal/application"
	"eda-in-golang/modules/depot/internal/grpc"
	"eda-in-golang/modules/depot/internal/handlers"
	"eda-in-golang/modules/depot/internal/logging"
	"eda-in-golang/modules/depot/internal/postgres"
	"eda-in-golang/modules/depot/internal/rest"
)

type Module struct{}

var _ monolith.Module = (*Module)(nil)

func (Module) Startup(ctx context.Context, srv monolith.Server) error {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	shoppingLists := postgres.NewShoppingListRepository("depot.shopping_lists", srv.DB())
	conn, err := grpc.Dial(ctx, srv.Config().Rpc.Address())
	if err != nil {
		return err
	}
	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(application.New(shoppingLists, stores, products, domainDispatcher),
		srv.Logger(),
	)
	orderHandlers := logging.LogDomainEventHandlerAccess(
		application.NewOrderHandlers(orders),
		srv.Logger(),
	)

	// setup Driver adapters
	if err := grpc.Register(ctx, app, srv.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, srv.Mux(), srv.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(srv.Mux()); err != nil {
		return err
	}
	handlers.RegisterOrderHandlers(orderHandlers, domainDispatcher)

	return nil
}
