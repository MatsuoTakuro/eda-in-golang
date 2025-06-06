package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"eda-in-golang/internal/config"
	"eda-in-golang/internal/logger"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/rpc"
	"eda-in-golang/internal/runner"
	"eda-in-golang/internal/web"
	"eda-in-golang/modules/baskets"
	"eda-in-golang/modules/cosec"
	"eda-in-golang/modules/customers"
	"eda-in-golang/modules/depot"
	"eda-in-golang/modules/notifications"
	"eda-in-golang/modules/ordering"
	"eda-in-golang/modules/payments"
	"eda-in-golang/modules/search"
	"eda-in-golang/modules/stores"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	// parse config/env/...
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	m := mono{cfg: cfg}

	// init infrastructure...
	// init db
	m.db, err = sql.Open("pgx", cfg.PG.Conn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(m.db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = m.db.PingContext(ctx); err != nil {
		return
	}
	// init nats & jetstream
	m.nc, err = nats.Connect(cfg.Nats.URL)
	if err != nil {
		return err
	}
	defer m.nc.Close()
	m.js, err = initJetStream(cfg.Nats, m.nc)
	if err != nil {
		return err
	}
	m.logger = initLogger(cfg)
	m.rpc = initRpc(cfg.Rpc)
	m.mux = initMux(cfg.Web)
	m.runner = runner.New(runner.CatchSignals())

	// init modules
	m.modules = []monolith.Module{
		&baskets.Module{},
		&customers.Module{},
		&depot.Module{},
		&notifications.Module{},
		&ordering.Module{},
		&payments.Module{},
		&stores.Module{},
		&search.Module{},
		&cosec.Module{},
	}

	if err = m.startupModules(); err != nil {
		return err
	}

	// Mount general web resources
	m.mux.Mount("/", http.FileServer(http.FS(web.WebUI)))

	fmt.Println("started mallbots application")
	defer fmt.Println("stopped mallbots application")

	m.runner.Add(
		m.runWeb,
		m.runRPC,
		m.runStream,
	)

	return m.runner.Run()
}

func initLogger(cfg config.AppConfig) zerolog.Logger {
	return logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    logger.Level(cfg.LogLevel),
	})
}

func initRpc(_ rpc.RpcConfig) *grpc.Server {
	srv := grpc.NewServer()
	reflection.Register(srv)

	return srv
}

func initMux(_ web.WebConfig) *chi.Mux {
	return chi.NewMux()
}

func initJetStream(cfg config.NatsConfig, nc *nats.Conn) (jetstream.JetStream, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     cfg.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", cfg.Stream)},
	})

	return js, err
}
