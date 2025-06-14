package cosec

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
	"eda-in-golang/internal/sec"
	"eda-in-golang/internal/tm"
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

func (Module) Startup(ctx context.Context, mono monolith.Server) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton("registry", func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := customerspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := depotpb.RegisterMessages(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	container.AddSingleton("logger", func(c di.Container) (any, error) {
		return mono.Logger(), nil
	})
	container.AddSingleton("stream", func(c di.Container) (any, error) {
		return jetstream.NewStream("cosec", mono.Config().Nats.Stream, mono.JS(), c.Get("logger").(zerolog.Logger)), nil
	})
	container.AddSingleton("db", func(c di.Container) (any, error) {
		return mono.DB(), nil
	})
	container.AddSingleton("outboxProcessor", func(c di.Container) (any, error) {
		return tm.NewOutboxProcessor(
			c.Get("stream").(am.RawMessageStream),
			pg.NewOutboxStore("cosec.outbox", c.Get("db").(*sql.DB)),
		), nil
	})
	container.AddScoped("tx", func(c di.Container) (any, error) {
		db := c.Get("db").(*sql.DB)
		return db.Begin()
	})
	container.AddScoped("txStream", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		outboxStore := pg.NewOutboxStore("cosec.outbox", tx)
		return am.WithRawMessageStreamMiddlewares(
			c.Get("stream").(am.RawMessageStream),
			tm.WithOutboxStream(outboxStore),
		), nil
	})
	container.AddScoped("eventStream", func(c di.Container) (any, error) {
		return am.NewEventStream(c.Get("registry").(registry.Registry), c.Get("txStream").(am.RawMessageStream)), nil
	})
	container.AddScoped("commandStream", func(c di.Container) (any, error) {
		return am.NewCommandStream(c.Get("registry").(registry.Registry), c.Get("txStream").(am.RawMessageStream)), nil
	})
	container.AddScoped("replyStream", func(c di.Container) (any, error) {
		return am.NewReplyStream(c.Get("registry").(registry.Registry), c.Get("txStream").(am.RawMessageStream)), nil
	})
	container.AddScoped("inboxMiddleware", func(c di.Container) (any, error) {
		tx := c.Get("tx").(*sql.Tx)
		inboxStore := pg.NewInboxStore("cosec.inbox", tx)
		return tm.WithInboxHandler(inboxStore), nil
	})
	container.AddScoped(di.SagaRepo, func(c di.Container) (any, error) {
		reg := c.Get("registry").(registry.Registry)
		return sec.NewRepository[*models.CreateOrderData](
			reg,
			pg.NewSagaStore(
				"cosec.sagas",
				c.Get("tx").(*sql.Tx),
				reg,
			),
		), nil
	})
	container.AddSingleton(di.Saga, func(c di.Container) (any, error) {
		return internal.NewCreateOrderSaga(), nil
	})

	// setup application
	container.AddScoped(di.Orchestrator, func(c di.Container) (any, error) {
		return logging.LogReplyHandlerAccess[*models.CreateOrderData](
			sec.NewOrchestrator[*models.CreateOrderData](
				c.Get(di.Saga).(sec.Saga[*models.CreateOrderData]),
				c.Get(di.SagaRepo).(sec.Repository[*models.CreateOrderData]),
				c.Get("commandStream").(am.CommandStream),
			),
			"CreateOrderSaga", mono.Logger(),
		), nil
	})
	container.AddScoped("integrationEventHandlers", func(c di.Container) (any, error) {
		return logging.LogEventHandlerAccess[ddd.Event](
			handlers.NewIntegrationEventHandlers(
				c.Get(di.Orchestrator).(sec.Orchestrator[*models.CreateOrderData]),
			),
			"IntegrationEvents", c.Get("logger").(zerolog.Logger),
		), nil
	})

	// setup Driver adapters
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	if err = handlers.RegisterReplyHandlersTx(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, container)

	return
}

func registrations(reg registry.Registry) (err error) {
	serde := registrar.NewJsonRegistrar(reg)

	// Saga data
	if err = serde.RegisterWithKey(internal.CreateOrderSagaName, models.CreateOrderData{}); err != nil {
		return err
	}

	return nil
}

func startOutboxProcessor(ctx context.Context, container di.Container) {
	outboxProcessor := container.Get("outboxProcessor").(tm.OutboxProcessor)
	logger := container.Get("logger").(zerolog.Logger)

	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("cosec outbox processor encountered an error")
		}
	}()
}
