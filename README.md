# gocron

Small Go cron wrapper around `github.com/robfig/cron` with context-aware jobs, optional timeouts, locking, and error handling hooks.

## Features
- Context-aware job execution.
- Optional per-job timeout.
- Pluggable lock interface to avoid concurrent runs.
- Error handler with execution stage information.
- Graceful shutdown that waits for running jobs.

## Install
```bash
go get github.com/anticrew/gocron
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

	gocron "github.com/anticrew/gocron"
)

func main() {
	l := slog.Default()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notifyCtx, cancelNotifyCtx := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancelNotifyCtx()

	c := gocron.NewCron(notifyCtx,
		gocron.WithTimeout(15*time.Second),
		gocron.WithDefaultHandler(gocron.SlogHandler(l, slog.LevelError)),
	)

	gocron.Must(c.Add("*/1 * * * * *", func(ctx context.Context) error {
		log.Println("@every 1s run")
		return nil
	})).WithName("1s ok")

	c.MustAdd("@every 1s", func(ctx context.Context) error {
		return errors.New("no data")
	}).WithName("1s err")

	c.MustAdd("@every 1m", func(ctx context.Context) error {
		time.Sleep(16 * time.Second)
		return ctx.Err()
	}).WithName("1s timeout")

	c.Start()

	<-notifyCtx.Done()

	if err := c.Shutdown(ctx); err != nil {
		l.Error("can't shutdown cron", slog.Any("error", err))
		return
	}
}
```

## Testing
See `ai-rules/test/SKILL.md` for unit test guidelines.
