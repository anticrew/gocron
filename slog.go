package gocron

import (
	"context"
	"log/slog"

	"github.com/anticrew/gocron/internal"
)

// SlogHandler logs job events and errors via slog.
type SlogHandler struct {
	log          *slog.Logger
	eventLeveler slog.Leveler
	errorLeveler slog.Leveler
}

// NewSlogHandler creates a slog-based handler with no event/error levels set by default.
// If specified log is nil, slog.Default will be used.
func NewSlogHandler(log *slog.Logger) *SlogHandler {
	return &SlogHandler{
		log: internal.WithDefault(log, slog.Default),
	}
}

// WithError sets the level used for error events. nil value disables error handling
func (s *SlogHandler) WithError(leveler slog.Leveler) *SlogHandler {
	s.errorLeveler = leveler
	return s
}

// WithEvent sets the level used for event events. nil value disables event handling
func (s *SlogHandler) WithEvent(leveler slog.Leveler) *SlogHandler {
	s.eventLeveler = leveler
	return s
}

// Handle logs a job event based on stage and error presence.
func (s *SlogHandler) Handle(event JobEvent) {
	if event.Error != nil {
		s.handleError(event)
		return
	}

	s.handleEvent(event)
}

func (s *SlogHandler) handleError(event JobEvent) {
	if s.errorLeveler == nil {
		return
	}

	var msg string

	switch event.Stage {
	case StageStart:
		msg = "can't start job"

	case StageExec:
		msg = "can't execute job"

	case StageFinish:
		msg = "can't finish job"
	}

	s.log.LogAttrs(context.Background(), s.errorLeveler.Level(), msg,
		slog.String("spec", event.JobSpec),
		slog.String("name", event.JobName),
		slog.Any("error", event.Error))
}

func (s *SlogHandler) handleEvent(event JobEvent) {
	if s.eventLeveler == nil {
		return
	}

	var msg string

	switch event.Stage {
	case StageStart:
		msg = "job started"

	case StageExec:
		msg = "job executed"

	case StageFinish:
		msg = "job finished"
	}

	s.log.LogAttrs(context.Background(), s.eventLeveler.Level(), msg,
		slog.String("spec", event.JobSpec),
		slog.String("name", event.JobName))
}
