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

	if err := paymentspb.RegisterMessages(reg); err != nil {
		return err
	}
	if err := orderingpb.RegisterMessages(reg); err != nil {
		return err
	}
	stream := jetstream.NewStream("payments", mono.Config().Nats.Stream, mono.JS(), mono.Logger())
	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)
	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()
	invoices := postgres.NewInvoiceRepository("payments.invoices", mono.DB())
	payments := postgres.NewPaymentRepository("payments.payments", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(invoices, payments, domainDispatcher),
		mono.Logger(),
	)
	domainEventHandlers := logging.LogEventHandlerAccess(
		handlers.NewDomainEventHandlers(eventStream),
		"DomainEvents", mono.Logger(),
	)
	integrationEventHandlers := logging.LogEventHandlerAccess(
		handlers.NewIntegrationHandlers(app),
		"IntegrationEvents", mono.Logger(),
	)
	commandHandlers := logging.LogCommandHandlerAccess(
		handlers.NewCommandHandlers(app),
		"Commands", mono.Logger(),
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
	if err := handlers.RegisterIntegrationEventHandlers(eventStream, integrationEventHandlers); err != nil {
		return err
	}
	handlers.SubscribeDomainEvents(domainDispatcher, domainEventHandlers)
	if err := handlers.SubscribeCommands(commandStream, commandHandlers); err != nil {
		return err
	}

	return nil
}
