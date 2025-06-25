package system

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"eda-in-golang/internal/config"
	"eda-in-golang/internal/logger"
	"eda-in-golang/internal/runner"
)

type system struct {
	cfg    config.AppConfig
	db     *sql.DB
	nc     *nats.Conn
	js     jetstream.JetStream
	mux    *chi.Mux
	rpc    *grpc.Server
	runner runner.Runner
	logger zerolog.Logger
}

var _ Service = (*system)(nil)

func NewSystem(cfg config.AppConfig) (*system, error) {
	s := &system{cfg: cfg}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initJS(); err != nil {
		return nil, err
	}

	s.initMux()
	s.initRpc()
	s.initRunner()
	s.initLogger()

	return s, nil
}

func (s *system) Config() config.AppConfig {
	return s.cfg
}

func (s *system) initDB() (err error) {
	s.db, err = sql.Open("pgx", s.cfg.PG.Conn)
	return err
}

func (s *system) MigrateDB(fs fs.FS) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(s.db, "."); err != nil {
		return err
	}
	return nil
}

func (s *system) DB() *sql.DB {
	return s.db
}

func (s *system) initJS() (err error) {
	s.nc, err = nats.Connect(s.cfg.Nats.URL)
	if err != nil {
		return err
	}
	s.js, err = jetstream.New(s.nc)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = s.js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     s.cfg.Nats.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", s.cfg.Nats.Stream)},
	})

	return err
}

func (s *system) JS() jetstream.JetStream {
	return s.js
}

func (s *system) initLogger() {
	s.logger = logger.New(logger.LogConfig{
		Environment: s.cfg.Environment,
		LogLevel:    logger.Level(s.cfg.LogLevel),
	})
}

func (s *system) Logger() zerolog.Logger {
	return s.logger
}

func (s *system) initMux() {
	s.mux = chi.NewMux()

	s.mux.Use(middleware.Heartbeat("/liveness"))
}

func (s *system) Mux() *chi.Mux {
	return s.mux
}

func (s *system) initRpc() {
	s.rpc = grpc.NewServer()
	reflection.Register(s.rpc)
}

func (s *system) RPC() *grpc.Server {
	return s.rpc
}

func (s *system) initRunner() {
	s.runner = runner.New(runner.CatchSignals())
}

func (s *system) Runner() runner.Runner {
	return s.runner
}

func (s *system) Nats() *nats.Conn {
	return s.nc
}

func (s *system) RunWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    s.cfg.Web.Address(),
		Handler: s.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("web server started; listening at http://localhost%s\n", s.cfg.Web.Port)
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}

func (s *system) RunRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("rpc server started")
		defer fmt.Println("rpc server shutdown")
		if err := s.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(s.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			// Force it to stop
			s.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (s *system) RunStream(ctx context.Context) error {
	closed := make(chan struct{})
	s.nc.SetClosedHandler(func(*nats.Conn) {
		close(closed)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("message stream started")
		defer fmt.Println("message stream stopped")
		<-closed
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		return s.nc.Drain()
	})
	return group.Wait()
}
