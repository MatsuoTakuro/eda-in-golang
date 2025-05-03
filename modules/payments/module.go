package payments

import (
	"context"

	"eda-in-golang/internal/monolith"
	"eda-in-golang/modules/payments/internal/application"
	"eda-in-golang/modules/payments/internal/grpc"
	"eda-in-golang/modules/payments/internal/logging"
	"eda-in-golang/modules/payments/internal/postgres"
	"eda-in-golang/modules/payments/internal/rest"
)

type Module struct{}

var _ monolith.Module = (*Module)(nil)

func (m Module) Startup(ctx context.Context, srv monolith.Server) error {
	// setup Driven adapters
	invoices := postgres.NewInvoiceRepository("payments.invoices", srv.DB())
	payments := postgres.NewPaymentRepository("payments.payments", srv.DB())
	conn, err := grpc.Dial(ctx, srv.Config().Rpc.Address())
	if err != nil {
		return err
	}
	orders := grpc.NewOrderRepository(conn)

	// setup application
	var app application.App
	app = application.New(invoices, payments, orders)
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
