package customers

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/modules/customers/internal/application"
	"eda-in-golang/modules/customers/internal/grpc"
	"eda-in-golang/modules/customers/internal/logging"
	"eda-in-golang/modules/customers/internal/postgres"
	"eda-in-golang/modules/customers/internal/rest"
)

type Module struct{}

var _ monolith.Module = (*Module)(nil)

func (m Module) Startup(ctx context.Context, srv monolith.Server) error {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	customers := postgres.NewCustomerRepository("customers.customers", srv.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers, domainDispatcher),
		srv.Logger(),
	)

	if err := grpc.RegisterServer(app, srv.RPC()); err != nil {
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
