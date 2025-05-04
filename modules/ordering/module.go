package ordering

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/modules/ordering/internal/application"
	"eda-in-golang/modules/ordering/internal/application/eventhandlers"
	"eda-in-golang/modules/ordering/internal/application/logging"
	infraEvtHdlrs "eda-in-golang/modules/ordering/internal/infra/eventhandlers"
	"eda-in-golang/modules/ordering/internal/infra/grpc"
	"eda-in-golang/modules/ordering/internal/infra/postgres"
	"eda-in-golang/modules/ordering/internal/rest"
)

type Module struct{}

var _ monolith.Module = (*Module)(nil)

func (Module) Startup(ctx context.Context, srv monolith.Server) error {
	// setup Driven (Outbound) adapters
	domainDispatcher := ddd.NewEventDispatcher()
	orderRepo := postgres.NewOrderRepository("ordering.orders", srv.DB())
	conn, err := grpc.Dial(ctx, srv.Config().Rpc.Address())
	if err != nil {
		return err
	}
	customerCli := grpc.NewCustomerClient(conn)
	paymentCli := grpc.NewPaymentClient(conn)
	invoiceCli := grpc.NewInvoiceClient(conn)
	shoppingCli := grpc.NewShoppingListClient(conn)
	notificationCli := grpc.NewNotificationClient(conn)

	// setup application with logging
	app := logging.NewUsecases(
		application.NewUsecases(orderRepo, customerCli, paymentCli, shoppingCli, domainDispatcher),
		srv.Logger(),
	)
	// setup application handlers with logging
	notificationEventHdlrs := logging.NewNotificationEventHandler(
		eventhandlers.NewNotification(notificationCli),
		srv.Logger(),
	)
	invoiceEventHdlrs := logging.NewInvoiceEventHandler(
		eventhandlers.NewInvoice(invoiceCli),
		srv.Logger(),
	)

	// setup Driver (Inbound) adapters
	if err := grpc.RegisterServer(app, srv.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, srv.Mux(), srv.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(srv.Mux()); err != nil {
		return err
	}
	infraEvtHdlrs.SubscribeForNotification(notificationEventHdlrs, domainDispatcher)
	infraEvtHdlrs.SubscribeForInvoice(invoiceEventHdlrs, domainDispatcher)

	return nil
}
