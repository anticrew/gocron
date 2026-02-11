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
1. Create base context
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

2. Create a cron instance with a parent context
```go
cron := gocron.NewCron(ctx)
```
or cron instance with optional timeout
```go
cron := gocron.NewCron(ctx, gocron.WithTimeout(15*time.Second))
```

3. Register jobs with schedules and, optionally, names or handlers.
```go
j, err := cron.Add("* * * * *", func(ctx context.Context) error {
	log.Println("tick")
	return nil
})
if err != nil {
	// handle error
}

j.WithName("#1")
```
```go
h := gocron.NewSlogHandler(slog.Default()).
    WithError(slog.LevelError).
    WithEvent(slog.LevelDebug)

cron.MustAdd("@every 1s", func(ctx context.Context) error {
	log.Println("tick")
	return nil
}).WithName("#2").WithHandler(h)
```
You can set default handler by passing `WithDefaultHandler` option to `NewCron` function. 

4. Start the scheduler and shut it down using the same context.
```go
cron.Start()

if err := cron.Shutdown(ctx); err != nil {
	slog.Default().Error("shutdown cron", slog.Any("error", err))
}
```
We recommend using a context with a timeout or deadline for `Shutdown` and ensuring it isn't already canceled.  
For a full example, e.g. signal-aware context, see `example` directory and `example/main.go`.

## Testing
See `ai-rules/test/SKILL.md` for unit test guidelines.
