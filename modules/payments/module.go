package payments

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/ordering/orderingpb"
	"eda-in-golang/modules/payments/internal/application"
	"eda-in-golang/modules/payments/internal/grpc"
	"eda-in-golang/modules/payments/internal/handlers"
	"eda-in-golang/modules/payments/internal/logging"
	"eda-in-golang/modules/payments/internal/postgres"
	"eda-in-golang/modules/payments/internal/rest"
	"eda-in-golang/modules/payments/paymentspb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Server) error {
	// setup Driven adapters
	reg := registry.New()

	if err := paymentspb.RegisterIntegrationEvents(reg); err != nil {
		return err
	}
	if err := orderingpb.RegisterIntegrationEvents(reg); err != nil {
		return err
	}
	eventStream := am.NewEventStream(reg, jetstream.NewStream("payments", mono.Config().Nats.Stream, mono.JS()))
	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()
	invoices := postgres.NewInvoiceRepository("payments.invoices", mono.DB())
	payments := postgres.NewPaymentRepository("payments.payments", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(invoices, payments, domainDispatcher),
		mono.Logger(),
	)
	orderHandlers := logging.LogEventHandlerAccess(
		application.NewOrderHandlers(app),
		"Order", mono.Logger(),
	)
	integrationEventHandlers := logging.LogEventHandlerAccess(
		application.NewIntegrationEventHandlers(eventStream),
		"IntegrationEvents", mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}
	if err := handlers.SubscribeOrderIntegrationEvents(orderHandlers, eventStream); err != nil {
		return err
	}
	handlers.SubscribeDomainEventsForIntegration(integrationEventHandlers, domainDispatcher)

	return nil
}
