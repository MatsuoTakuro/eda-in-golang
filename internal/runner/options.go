package runner

import (
	"context"
)

type RunnerOption func(c *runnerCfg)

func ParentContext(ctx context.Context) RunnerOption {
	return func(c *runnerCfg) {
		c.parentCtx = ctx
	}
}

func CatchSignals() RunnerOption {
	return func(c *runnerCfg) {
		c.catchSignals = true
	}
}
