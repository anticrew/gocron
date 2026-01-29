package gocron

import (
	"context"
	"time"
)

type Cmd func(ctx context.Context) error

type Cron interface {
	Add(spec string, cmd Cmd) Job
	Start() error
	Shutdown(ctx context.Context) error
}

type Job interface {
	WithName(name string) Job
	WithTimeout(t time.Duration) Job
	WithLock(lock Lock) Job
	WithHandler(h Handler) Job
}

type Lock interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

type Stage int8

const (
	StageStart Stage = 1 << iota
	StageRun
	StageFinish
)

type Handler func(spec, name string, state Stage, err error)
