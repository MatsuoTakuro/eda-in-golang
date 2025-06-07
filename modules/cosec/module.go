package cosec

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
	"eda-in-golang/internal/sec"
	"eda-in-golang/modules/cosec/internal"
	"eda-in-golang/modules/cosec/internal/handlers"
	"eda-in-golang/modules/cosec/internal/logging"
	"eda-in-golang/modules/cosec/internal/models"
	"eda-in-golang/modules/customers/customerspb"
	"eda-in-golang/modules/depot/depotpb"
	"eda-in-golang/modules/ordering/orderingpb"
	"eda-in-golang/modules/payments/paymentspb"
)

type Module struct{}

func (Module) Startup(ctx context.Context, srv monolith.Server) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = registerMessages(reg); err != nil {
		return err
	}
	if err = orderingpb.RegisterMessages(reg); err != nil {
		return err
	}
	if err = customerspb.RegisterMessages(reg); err != nil {
		return err
	}
	if err = depotpb.RegisterMessages(reg); err != nil {
		return err
	}
	if err = paymentspb.RegisterMessages(reg); err != nil {
		return err
	}
	stream := jetstream.NewStream("cosec", srv.Config().Nats.Stream, srv.JS(), srv.Logger())
	eventStream := am.NewEventStream(reg, stream)
	commandStream := am.NewCommandStream(reg, stream)
	replyStream := am.NewReplyStream(reg, stream)
	sagaStore := pg.NewSagaStore("cosec.sagas", srv.DB(), reg)
	sagaRepo := sec.NewRepository[*models.CreateOrderData](reg, sagaStore)

	// setup application
	orchestrator := logging.LogReplyHandlerAccess(
		sec.NewOrchestrator(internal.NewCreateOrderSaga(), sagaRepo, commandStream),
		"CreateOrderSaga", srv.Logger(),
	)
	integrationEventHandlers := logging.LogEventHandlerAccess(
		handlers.NewIntegrationEventHandlers(orchestrator),
		"IntegrationEvents", srv.Logger(),
	)

	// setup Driver adapters
	if err = handlers.SubscribeIntegrationEvents(eventStream, integrationEventHandlers); err != nil {
		return err
	}
	if err = handlers.SubscribeReplies(replyStream, orchestrator); err != nil {
		return err
	}

	return
}

func registerMessages(reg registry.Registry) (err error) {
	regtr := registrar.NewJsonRegistrar(reg)

	// Saga data
	if err = regtr.RegisterWithKey(internal.CreateOrderSagaName, models.CreateOrderData{}); err != nil {
		return err
	}

	return nil
}
