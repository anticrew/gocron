package gocron

import "errors"

var (
	ErrCommandIsNil   = errors.New("command is nil")
	ErrCronNotRunning = errors.New("cron is not running")
)
