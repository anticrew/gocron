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
