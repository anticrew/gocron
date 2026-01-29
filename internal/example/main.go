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
		//log.Println("@every 1s run")
		return nil
	}).WithName("1h ok")

	c.Add("@every 1s", func(ctx context.Context) error {
		return errors.New("no data")
	}).WithName("1s err")

	c.Add("@every 1m", func(ctx context.Context) error {
		time.Sleep(16 * time.Second)
		return ctx.Err()
	}).WithName("1s timeout")

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	<-notifyCtx.Done()

	if err := c.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
