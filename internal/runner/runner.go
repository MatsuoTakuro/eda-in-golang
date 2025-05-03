package runner

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type RunFunc func(ctx context.Context) error

type Runner interface {
	Add(fns ...RunFunc)
	Run() error
	Context() context.Context
	CancelFunc() context.CancelFunc
}

type runner struct {
	ctx    context.Context
	fns    []RunFunc
	cancel context.CancelFunc
}

type runnerCfg struct {
	parentCtx    context.Context
	catchSignals bool
}

func New(options ...RunnerOption) Runner {
	cfg := &runnerCfg{
		parentCtx:    context.Background(),
		catchSignals: false,
	}

	for _, option := range options {
		option(cfg)
	}

	w := &runner{
		fns: []RunFunc{},
	}
	w.ctx, w.cancel = context.WithCancel(cfg.parentCtx)
	if cfg.catchSignals {
		w.ctx, w.cancel = signal.NotifyContext(w.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return w
}

func (w *runner) Add(fns ...RunFunc) {
	w.fns = append(w.fns, fns...)
}

func (w runner) Run() (err error) {
	g, ctx := errgroup.WithContext(w.ctx)

	g.Go(func() error {
		<-ctx.Done()
		w.cancel()
		return nil
	})

	for _, fn := range w.fns {
		fn := fn
		g.Go(func() error { return fn(ctx) })
	}

	return g.Wait()
}

func (w runner) Context() context.Context {
	return w.ctx
}

func (w runner) CancelFunc() context.CancelFunc {
	return w.cancel
}
