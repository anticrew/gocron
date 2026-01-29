# go-cron

Small Go cron wrapper around `github.com/robfig/cron` with context-aware jobs, optional timeouts, locking, and error handling hooks.

## Features
- Context-aware job execution.
- Optional per-job timeout.
- Pluggable lock interface to avoid concurrent runs.
- Error handler with execution stage information.
- Graceful shutdown that waits for running jobs.

## Install
```bash
go get github.com/KirillBogatikov/go-cron
```

## Quick start
```go
package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	gocron "github.com/KirillBogatikov/go-cron"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notifyCtx, cancelNotifyCtx := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancelNotifyCtx()

	c := gocron.NewCron(notifyCtx,
		gocron.WithTimeout(15*time.Second),
		gocron.WithDefaultHandler(gocron.SlogHandler(slog.Default(), slog.LevelError)),
	)

	c.Add("@every 1s", func(ctx context.Context) error {
		return nil
	}).WithName("fast-ok")

	c.Add("@every 1s", func(ctx context.Context) error {
		return errors.New("no data")
	}).WithName("fast-error")

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	<-notifyCtx.Done()

	if err := c.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
```

## API overview
### Cron
```go
type Cron interface {
	Add(spec string, cmd Cmd) Job
	Start() error
	Shutdown(ctx context.Context) error
}
```

### Job options
```go
type Job interface {
	WithName(name string) Job
	WithTimeout(t time.Duration) Job
	WithLock(lock Lock) Job
	WithHandler(h Handler) Job
}
```

### Handlers and stages
- Use `WithDefaultHandler` to set a default job handler for the cron.
- Use `WithHandler` to override per job.
- Handler receives the stage: `StageStart`, `StageRun`, `StageFinish`.
- `SlogHandler` is a ready-to-use handler based on `log/slog`.

### Locking
Provide a custom lock implementation:
```go
type Lock interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}
```

## Testing
See `TESTS.md` for unit test guidelines.
