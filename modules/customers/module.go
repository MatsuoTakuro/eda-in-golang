package customers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/customers/internal/application"
	"eda-in-golang/modules/customers/internal/grpc"
	"eda-in-golang/modules/customers/internal/handlers"
	"eda-in-golang/modules/customers/internal/logging"
	"eda-in-golang/modules/customers/internal/postgres"
	"eda-in-golang/modules/customers/internal/rest"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Server) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = customerspb.RegisterMessages(reg); err != nil {
		return err
	}
	stream := jetstream.NewStream("customers", mono.Config().Nats.Stream, mono.JS())
	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	customers := postgres.NewCustomerRepository("customers.customers", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers, domainDispatcher),
		mono.Logger(),
	)
	domainEventHandlers := logging.LogEventHandlerAccess(
		application.NewDomainEventHandlers(eventStream),
		"DomainEvents", mono.Logger(),
	)
	commandHandlers := logging.LogCommandHandlerAccess(
		application.NewCommandHandlers(app),
		"Commands", mono.Logger(),
	)
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}
	handlers.SubscribeDomainEvents(domainEventHandlers, domainDispatcher)
	if err = handlers.SubscribeCommands(commandStream, commandHandlers); err != nil {
		return err
	}
	return nil
}
