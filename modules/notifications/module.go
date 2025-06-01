package notifications

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/registry"
	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/notifications/internal/application"
	"eda-in-golang/modules/notifications/internal/grpc"
	"eda-in-golang/modules/notifications/internal/handlers"
	"eda-in-golang/modules/notifications/internal/logging"
	"eda-in-golang/modules/notifications/internal/postgres"
	"eda-in-golang/modules/ordering/orderingpb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Server) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = customerspb.RegisterMessages(reg); err != nil {
		return err
	}
	if err = orderingpb.RegisterMessages(reg); err != nil {
		return err
	}
	eventStream := am.NewEventStream(reg, jetstream.NewStream("notifications", mono.Config().Nats.Stream, mono.JS(), mono.Logger()))
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	customers := postgres.NewCustomerCacheRepository("notifications.customers_cache", mono.DB(), grpc.NewCustomerRepository(conn))

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers),
		mono.Logger(),
	)
	customerHandlers := logging.LogEventHandlerAccess(
		application.NewCustomerHandlers(customers),
		"Customer", mono.Logger(),
	)
	orderHandlers := logging.LogEventHandlerAccess(
		application.NewOrderHandlers(app),
		"Order", mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err = handlers.SubscribeCustomerIntegrationEvents(customerHandlers, eventStream); err != nil {
		return err
	}
	if err = handlers.SubscribeOrderIntegrationEvents(orderHandlers, eventStream); err != nil {
		return err
	}

	return nil
}
