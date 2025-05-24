package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"eda-in-golang/internal/config"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/runner"
)

type mono struct {
	cfg     config.AppConfig
	db      *sql.DB
	nc      *nats.Conn
	js      jetstream.JetStream
	logger  zerolog.Logger
	modules []monolith.Module
	mux     *chi.Mux
	rpc     *grpc.Server
	runner  runner.Runner
}

var _ monolith.Server = (*mono)(nil)

func (m *mono) Config() config.AppConfig {
	return m.cfg
}

func (m *mono) DB() *sql.DB {
	return m.db
}

func (m *mono) JS() jetstream.JetStream {
	return m.js
}

func (m *mono) Logger() zerolog.Logger {
	return m.logger
}

func (m *mono) Mux() *chi.Mux {
	return m.mux
}

func (m *mono) RPC() *grpc.Server {
	return m.rpc
}

func (m *mono) Runner() runner.Runner {
	return m.runner
}

func (m *mono) startupModules() error {
	for _, module := range m.modules {
		if err := module.Startup(m.Runner().Context(), m); err != nil {
			return err
		}
	}

	return nil
}

func (m *mono) runWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    m.cfg.Web.Address(),
		Handler: m.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("web server started")
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), m.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

func (m *mono) runRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", m.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("rpc server started")
		defer fmt.Println("rpc server shutdown")
		if err := m.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			m.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(m.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			// Force it to stop
			m.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (m *mono) runStream(ctx context.Context) error {
	closed := make(chan struct{})
	m.nc.SetClosedHandler(func(*nats.Conn) { // when the connection is closed (3)
		close(closed) // close the channel (4)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("message stream started")       // start the message stream (1)
		defer fmt.Println("message stream stopped") // stop the message stream (6)
		<-closed                                    // wait for the connection to be closed (5)
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done() // wait for the context to be done or cancelled (2)
		return m.nc.Drain()
	})
	return group.Wait() // wait for the group to finish (7)
}
