package gocron

import (
	"log/slog"
)

func SlogHandler(log *slog.Logger, level slog.Level) Handler {
	return func(spec, name string, stage Stage, err error) {
		var msg string

		switch stage {
		case StageStart:
			msg = "can't start job"

		case StageRun:
			msg = "can't run job"

		case StageFinish:
			msg = "can't finish job"
		}

		log.LogAttrs(nil, level, msg,
			slog.String("spec", spec),
			slog.String("name", name),
			slog.Any("error", err))
	}
}
