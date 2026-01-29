package gocron

import (
	"context"
	"errors"
	"sync"
	"time"

	c "github.com/robfig/cron"
)

type defaults struct {
	handler Handler
	timeout time.Duration
}

type cron struct {
	baseCtx context.Context
	cron    *c.Cron
	errors  []error

	defaults defaults
	wg       *sync.WaitGroup
}

type Option func(*cron)

func WithDefaultHandler(h Handler) Option {
	return func(c *cron) {
		c.defaults.handler = h
	}
}

func WithTimeout(t time.Duration) Option {
	return func(c *cron) {
		c.defaults.timeout = t
	}
}

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

func (c *cron) Add(spec string, cmd Cmd) Job {
	j := newJob(spec, c.baseCtx, cmd)
	j.WithHandler(c.defaults.handler)
	j.WithTimeout(c.defaults.timeout)
	j.withWaitGroup(c.wg)

	if err := c.cron.AddJob(spec, j); err != nil {
		c.errors = append(c.errors)
	}

	return j
}

func (c *cron) Start() error {
	if len(c.errors) > 0 {
		return errors.Join(c.errors...)
	}

	c.cron.Start()
	return nil
}

func (c *cron) Shutdown(ctx context.Context) error {
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
