package gocron

import (
	"context"
	"log/slog"
)

// SlogHandler logs job events and errors via slog.
type SlogHandler struct {
	log          *slog.Logger
	traceLeveler slog.Leveler
	errorLeveler slog.Leveler
}

// NewSlogHandler creates a slog-based handler with no trace/error levels set.
func NewSlogHandler(log *slog.Logger) *SlogHandler {
	return &SlogHandler{
		log: log,
	}
}

// WithError sets the level used for error events.
func (s *SlogHandler) WithError(leveler slog.Leveler) *SlogHandler {
	s.errorLeveler = leveler
	return s
}

// WithTrace sets the level used for trace events.
func (s *SlogHandler) WithTrace(leveler slog.Leveler) *SlogHandler {
	s.traceLeveler = leveler
	return s
}

// Handle logs a job event based on stage and error presence.
func (s *SlogHandler) Handle(spec, name string, stage Stage, err error) {
	if err != nil {
		s.handleError(spec, name, stage, err)
		return
	}

	s.handleTrace(spec, name, stage)
}

func (s *SlogHandler) handleError(spec, name string, stage Stage, err error) {
	if s.errorLeveler == nil {
		return
	}

	var msg string

	switch stage {
	case StageStart:
		msg = "can't start job"

	case StageRun:
		msg = "can't run job"

	case StageFinish:
		msg = "can't finish job"
	}

	s.log.LogAttrs(context.Background(), s.errorLeveler.Level(), msg,
		slog.String("spec", spec),
		slog.String("name", name),
		slog.Any("error", err))
}

func (s *SlogHandler) handleTrace(spec, name string, stage Stage) {
	if s.traceLeveler == nil {
		return
	}

	var msg string

	switch stage {
	case StageStart:
		msg = "job started"

	case StageRun:
		msg = "job executed"

	case StageFinish:
		msg = "job finished"
	}

	s.log.LogAttrs(context.Background(), s.traceLeveler.Level(), msg,
		slog.String("spec", spec),
		slog.String("name", name))
}
