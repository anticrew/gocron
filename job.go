package gocron

import (
	"context"
	"sync"
	"time"

	"github.com/KirillBogatikov/go-cron/internal"
)

type job struct {
	spec, name string

	baseCtx    context.Context
	newContext internal.ContextFactory

	wg   *sync.WaitGroup
	lock Lock

	cmd     Cmd
	handler Handler
}

func newJob(baseCtx context.Context, spec string, cmd Cmd) *job {
	return &job{
		spec:    spec,
		name:    internal.RandName(8),
		baseCtx: baseCtx,
		cmd:     cmd,
	}
}

// Run executes the job command with lock and handler hooks.
// Exported for compliance with github.com/robfig/cron's Job interface and shouldn't be called manually
func (j *job) Run() {
	j.wg.Add(1)
	defer j.wg.Done()

	ctx, cancel := context.WithCancel(j.baseCtx)
	defer cancel()

	if !j.acquireLock(ctx) {
		return
	}

	defer j.releaseLock(ctx)

	cmdCtx, cancelCmdCtx := j.newContext(ctx)
	defer cancelCmdCtx()

	if err := j.cmd(cmdCtx); err != nil {
		j.handle(StageRun, err)
	}
}

// WithTimeout sets the job timeout; non-positive value disables timeout
func (j *job) WithTimeout(t time.Duration) Job {
	var f internal.ContextFactory

	if t > 0 {
		f = internal.TimeoutContextFactory(t)
	} else {
		f = internal.CancelContextFactory()
	}

	j.newContext = f
	return j
}

// WithLock sets the lock used to guard concurrent runs
func (j *job) WithLock(lock Lock) Job {
	j.lock = lock
	return j
}

// WithHandler sets the error handler used by this job; nil disabled error handling
func (j *job) WithHandler(h Handler) Job {
	j.handler = h
	return j
}

// WithName sets the human-readable name used in handlers
func (j *job) WithName(name string) Job {
	j.name = name
	return j
}

func (j *job) withWaitGroup(wg *sync.WaitGroup) {
	j.wg = wg
}

func (j *job) acquireLock(ctx context.Context) bool {
	if j.lock == nil {
		return true
	}

	if err := j.lock.Lock(ctx); err != nil {
		j.handle(StageStart, err)
		return false
	}

	return true
}

func (j *job) releaseLock(ctx context.Context) {
	if j.lock == nil {
		return
	}

	if err := j.lock.Unlock(ctx); err != nil {
		j.handle(StageFinish, err)
	}
}

func (j *job) handle(stage Stage, err error) {
	if j.handler == nil {
		return
	}

	j.handler(j.spec, j.name, stage, err)
}
