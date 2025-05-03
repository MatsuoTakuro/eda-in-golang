package stores

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/modules/stores/internal/application"
	"eda-in-golang/modules/stores/internal/grpc"
	"eda-in-golang/modules/stores/internal/logging"
	"eda-in-golang/modules/stores/internal/postgres"
	"eda-in-golang/modules/stores/internal/rest"
)

type Module struct{}

var _ monolith.Module = (*Module)(nil)

func (m *Module) Startup(ctx context.Context, srv monolith.Server) error {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	stores := postgres.NewStoreRepository("stores.stores", srv.DB())
	participatingStores := postgres.NewParticipatingStoreRepository("stores.stores", srv.DB())
	products := postgres.NewProductRepository("stores.products", srv.DB())

	// setup application
	var app application.App
	app = application.New(stores, participatingStores, products, domainDispatcher)
	app = logging.LogApplicationAccess(app, srv.Logger())

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, srv.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, srv.Mux(), srv.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(srv.Mux()); err != nil {
		return err
	}

	return nil
}
