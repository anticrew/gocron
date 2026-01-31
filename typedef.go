package gocron

import (
	"context"
	"time"
)

// Cmd is the signature for a cron job command
type Cmd func(ctx context.Context) error

// Cron schedules and runs jobs
type Cron interface {
	// Add registers a job with the given cron spec.
	// Look at github.com/robfig/cron documentation for details about spec format
	Add(spec string, cmd Cmd) (Job, error)

	// MustAdd registers a job with the given cron spec like Add, but panics on any error.
	// Look at github.com/robfig/cron documentation for details about spec format
	MustAdd(spec string, cmd Cmd) Job

	// Start begins scheduling jobs.
	// It should be called once, next calls without call Shutdown before will be ignored
	Start()

	// Shutdown stops scheduling and waits for running jobs to finish or context cancellation.
	// It should be called once, next calls without call Start before will be ignored
	Shutdown(ctx context.Context) error
}

// Job configures a scheduled job
type Job interface {
	// WithName sets the human-readable name used in handlers
	WithName(name string) Job
	// WithTimeout sets the job timeout; non-positive value disables timeout
	WithTimeout(t time.Duration) Job
	// WithLock sets the lock used to guard concurrent runs
	WithLock(lock Lock) Job
	// WithHandler sets the error handler used by this job; nil disabled error handling
	WithHandler(h Handler) Job
}

// Lock guards concurrent job runs
type Lock interface {
	// Lock acquires the lock
	Lock(ctx context.Context) error
	// Unlock releases the lock
	Unlock(ctx context.Context) error
}

// Stage identifies the job lifecycle step for handler callbacks
type Stage int8

const (
	// StageStart indicates lock and start stage
	StageStart Stage = 1 << iota
	// StageExec indicates job execution stage
	StageExec
	// StageFinish indicates unlock and finish stage
	StageFinish
)

type JobEvent struct {
	JobSpec string
	JobName string
	Stage   Stage
	Error   error
}

// Handler receives job events and errors
type Handler interface {
	Handle(event JobEvent)
}

// HandlerFunc adapts a function to a Handler
type HandlerFunc func(event JobEvent)

// Handle calls the wrapped function
func (f HandlerFunc) Handle(event JobEvent) {
	f(event)
}
