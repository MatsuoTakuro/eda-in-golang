package notifications

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/system"
	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/notifications/internal/application"
	"eda-in-golang/modules/notifications/internal/grpc"
	"eda-in-golang/modules/notifications/internal/handlers"
	"eda-in-golang/modules/notifications/internal/logging"
	"eda-in-golang/modules/notifications/internal/postgres"
	"eda-in-golang/modules/ordering/orderingpb"
)

func Root(ctx context.Context, svc system.Service) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = customerspb.Registrations(reg); err != nil {
		return err
	}
	if err = orderingpb.Registrations(reg); err != nil {
		return err
	}
	eventStream := am.NewEventStream(reg,
		jetstream.NewStream("notifications", svc.Config().Nats.Stream, svc.JS(), svc.Logger()),
	)
	conn, err := grpc.Dial(ctx, svc.Config().Rpc.Service("CUSTOMERS"))
	if err != nil {
		return err
	}
	customers := postgres.NewCustomerCacheRepository("notifications.customers_cache", svc.DB(), grpc.NewCustomerRepository(conn))

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers),
		svc.Logger(),
	)
	integrationEventHandlers := logging.LogEventHandlerAccess[ddd.Event](
		handlers.NewIntegrationEventHandlers(app, customers),
		"IntegrationEvents", svc.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, svc.RPC()); err != nil {
		return err
	}
	if err = handlers.RegisterIntegrationEventHandlers(eventStream, integrationEventHandlers); err != nil {
		return err
	}

	return nil
}
