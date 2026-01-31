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
1. Create a cron instance with a parent context, optional timeout, and default error handler.
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

cron := gocron.NewCron(ctx,
	gocron.WithTimeout(15*time.Second),
	gocron.WithDefaultHandler(gocron.SlogHandler(slog.Default(), slog.LevelError)),
)
```

2. Register jobs with schedules and names.
```go
gocron.Must(cron.Add("*/1 * * * * *", func(ctx context.Context) error {
	log.Println("tick")
	return nil
})).WithName("1s ok")
```

3. Start the scheduler and shut it down using the same context.
```go
cron.Start()

if err := cron.Shutdown(ctx); err != nil {
	slog.Default().Error("shutdown cron", slog.Any("error", err))
}
```

For a full, signal-aware example, see `internal/example` and `internal/example/main.go`.

## Testing
See `ai-rules/test/SKILL.md` for unit test guidelines.
