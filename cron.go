package gocron

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	c "github.com/robfig/cron"
)

type defaults struct {
	handler Handler
	timeout time.Duration
}

type cron struct {
	started atomic.Bool

	baseCtx context.Context
	cron    *c.Cron

	defaults defaults
	wg       *sync.WaitGroup
}

type Option func(*cron)

// WithDefaultHandler sets the default handler used by all jobs.
// This handler can be overwritten by Job.WithHandler
func WithDefaultHandler(h Handler) Option {
	return func(c *cron) {
		c.defaults.handler = h
	}
}

// WithTimeout sets the default timeout used by all jobs.
// This timeout can be overwritten by Job.WithTimeout
func WithTimeout(t time.Duration) Option {
	return func(c *cron) {
		c.defaults.timeout = t
	}
}

// NewCron creates a cron with the provided context and options
func NewCron(ctx context.Context, options ...Option) Cron {
	cr := &cron{
		baseCtx: ctx,
		cron:    c.New(),
		wg:      &sync.WaitGroup{},
	}

	for _, option := range options {
		option(cr)
	}

	return cr
}

// Add registers a job with the given cron spec.
// Look at github.com/robfig/cron documentation for details about spec format
func (c *cron) Add(spec string, cmd Cmd) (Job, error) {
	j := newJob(c.baseCtx, spec, cmd)
	j.WithHandler(c.defaults.handler)
	j.WithTimeout(c.defaults.timeout)
	j.withWaitGroup(c.wg)

	if err := c.cron.AddJob(spec, j); err != nil {
		return nil, err
	}

	return j, nil
}

func (c *cron) MustAdd(spec string, cmd Cmd) Job {
	return Must(c.Add(spec, cmd))
}

// Start begins scheduling jobs.
// It should be called once, next calls without call Shutdown before will be ignored
func (c *cron) Start() {
	if c.started.Swap(true) {
		return
	}

	c.cron.Start()
}

// Shutdown stops scheduling and waits for running jobs to finish or context cancellation.
// It should be called once, next calls without call Start before will be ignored
func (c *cron) Shutdown(ctx context.Context) error {
	if !c.started.Swap(false) {
		return nil
	}

	c.cron.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.wait():
		return nil
	}
}

func (c *cron) wait() <-chan struct{} {
	ch := make(chan struct{}, 1)

	go func() {
		c.wg.Wait()
		close(ch)
	}()

	return ch
}
