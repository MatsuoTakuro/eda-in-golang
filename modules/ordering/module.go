package ordering

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
	"eda-in-golang/modules/baskets/basketspb"
	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/ordering/internal/application"
	"eda-in-golang/modules/ordering/internal/application/logging"
	"eda-in-golang/modules/ordering/internal/application/msghandlers"
	"eda-in-golang/modules/ordering/internal/domain"
	"eda-in-golang/modules/ordering/internal/infra/grpc"
	"eda-in-golang/modules/ordering/internal/infra/rest"
	"eda-in-golang/modules/ordering/orderingpb"
)

type Module struct{}

var _ monolith.Module = (*Module)(nil)

func (Module) Startup(ctx context.Context, srv monolith.Server) error {
	// setup Driven (Outbound) adapters
	reg := registry.New()
	if err := registerDomainEvents(reg); err != nil {
		return err
	}
	if err := basketspb.RegisterMessages(reg); err != nil {
		return err
	}
	if err := orderingpb.RegisterMessages(reg); err != nil {
		return err
	}
	if err := depotpb.RegisterMessages(reg); err != nil {
		return err
	}
	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()
	stream := jetstream.NewStream("ordering", srv.Config().Nats.Stream, srv.JS())
	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)
	aggregateStore := es.AggregateStoreWithMiddleware(
		pg.NewEventStore("ordering.events", srv.DB(), reg),
		pg.WithSnapshotStore("ordering.snapshots", srv.DB(), reg),
	)
	orderRepo := es.NewAggregateRepository[*domain.Order](
		domain.OrderAggregate,
		reg,
		aggregateStore,
	)

	conn, err := grpc.Dial(ctx, srv.Config().Rpc.Address())
	if err != nil {
		return err
	}
	shoppingCli := grpc.NewShoppingListClient(conn)

	// setup application with logging
	us := logging.NewUsecases(
		application.NewUsecases(orderRepo, shoppingCli, domainDispatcher),
		srv.Logger(),
	)
	domainEvtHdlr := logging.NewEventHandler(
		msghandlers.NewDomainEvents(eventStream),
		"DomainEvents", srv.Logger(),
	)
	integrationEvtHdlr := logging.NewEventHandler(
		msghandlers.NewIntegrationEvents(us),
		"IntegrationEvents", srv.Logger(),
	)
	cmdHdlr := logging.NewCommandHandler(
		msghandlers.NewCommands(us),
		"Commands", srv.Logger(),
	)

	// setup Driver (Inbound) adapters
	if err := grpc.RegisterServer(us, srv.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, srv.Mux(), srv.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(srv.Mux()); err != nil {
		return err
	}
	msghandlers.SubscribeDomainEvents(domainDispatcher, domainEvtHdlr)
	if err = msghandlers.SubscribeIntegrationEvents(eventStream, integrationEvtHdlr); err != nil {
		return err
	}
	if err = msghandlers.SubscribeCommands(commandStream, cmdHdlr); err != nil {
		return err
	}
	return nil
}

func registerDomainEvents(reg registry.Registry) error {
	regtr := registrar.NewJsonRegistrar(reg)

	// Order
	if err := regtr.Register(domain.Order{}, func(v any) error {
		order := v.(*domain.Order)
		order.Aggregate = es.NewAggregate("", domain.OrderAggregate)
		return nil
	}); err != nil {
		return err
	}
	// order events
	if err := regtr.Register(domain.OrderCreated{}); err != nil {
		return err
	}
	if err := regtr.Register(domain.OrderCanceled{}); err != nil {
		return err
	}
	if err := regtr.Register(domain.OrderReadied{}); err != nil {
		return err
	}
	if err := regtr.Register(domain.OrderCompleted{}); err != nil {
		return err
	}
	// order snapshots
	if err := regtr.Register(domain.OrderV1{}); err != nil {
		return err
	}

	return nil
}
