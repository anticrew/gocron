package gocron

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

type slogRecord struct {
	level slog.Level
	msg   string
	attrs map[string]any
}

type slogCaptureHandler struct {
	log []slogRecord
}

func (h *slogCaptureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *slogCaptureHandler) Handle(_ context.Context, r slog.Record) error {
	record := slogRecord{
		level: r.Level,
		msg:   r.Message,
		attrs: map[string]any{},
	}
	r.Attrs(func(attr slog.Attr) bool {
		record.attrs[attr.Key] = attr.Value.Any()
		return true
	})

	h.log = append(h.log, record)
	return nil
}

func (h *slogCaptureHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *slogCaptureHandler) WithGroup(string) slog.Handler {
	return h
}

func (h *slogCaptureHandler) assert(t *testing.T, records []slogRecord) {
	assert.Equal(t, records, h.log)
}

func TestSlogHandlerHandle(t *testing.T) {
	t.Parallel()

	type levelers struct {
		error slog.Leveler
		event slog.Leveler
	}

	tests := []struct {
		name     string
		event    JobEvent
		levelers levelers
		expected []slogRecord
	}{
		{
			name: "logs error for start stage",
			event: JobEvent{
				JobSpec: "@every 1s",
				JobName: "cleanup",
				Stage:   StageStart,
				Error:   assert.AnError,
			},
			levelers: levelers{
				error: slog.LevelError,
			},
			expected: []slogRecord{
				{
					level: slog.LevelError,
					msg:   "can't start job",
					attrs: map[string]any{
						"spec":  "@every 1s",
						"name":  "cleanup",
						"error": assert.AnError,
					},
				},
			},
		},
		{
			name: "logs error for exec stage",
			event: JobEvent{
				JobSpec: "@every 1s",
				JobName: "cleanup",
				Stage:   StageExec,
				Error:   assert.AnError,
			},
			levelers: levelers{
				error: slog.LevelError,
			},
			expected: []slogRecord{
				{
					level: slog.LevelError,
					msg:   "can't execute job",
					attrs: map[string]any{
						"spec":  "@every 1s",
						"name":  "cleanup",
						"error": assert.AnError,
					},
				},
			},
		},
		{
			name: "logs error for finish stage",
			event: JobEvent{
				JobSpec: "@every 1s",
				JobName: "cleanup",
				Stage:   StageFinish,
				Error:   assert.AnError,
			},
			levelers: levelers{
				error: slog.LevelError,
			},
			expected: []slogRecord{
				{
					level: slog.LevelError,
					msg:   "can't finish job",
					attrs: map[string]any{
						"spec":  "@every 1s",
						"name":  "cleanup",
						"error": assert.AnError,
					},
				},
			},
		},
		{
			name: "logs event for start stage",
			event: JobEvent{
				JobSpec: "0 0 * * *",
				JobName: "daily",
				Stage:   StageStart,
			},
			levelers: levelers{
				event: slog.LevelInfo,
			},
			expected: []slogRecord{
				{
					level: slog.LevelInfo,
					msg:   "job started",
					attrs: map[string]any{
						"spec": "0 0 * * *",
						"name": "daily",
					},
				},
			},
		},
		{
			name: "logs event for exec stage",
			event: JobEvent{
				JobSpec: "0 0 * * *",
				JobName: "daily",
				Stage:   StageExec,
			},
			levelers: levelers{
				event: slog.LevelInfo,
			},
			expected: []slogRecord{
				{
					level: slog.LevelInfo,
					msg:   "job executed",
					attrs: map[string]any{
						"spec": "0 0 * * *",
						"name": "daily",
					},
				},
			},
		},
		{
			name: "logs event for finish stage",
			event: JobEvent{
				JobSpec: "0 0 * * *",
				JobName: "daily",
				Stage:   StageFinish,
			},
			levelers: levelers{
				event: slog.LevelInfo,
			},
			expected: []slogRecord{
				{
					level: slog.LevelInfo,
					msg:   "job finished",
					attrs: map[string]any{
						"spec": "0 0 * * *",
						"name": "daily",
					},
				},
			},
		},
		{
			name: "skips logging when error level is nil",
			event: JobEvent{
				JobSpec: "@hourly",
				JobName: "backup",
				Stage:   StageFinish,
				Error:   assert.AnError,
			},
			expected: nil,
		},
		{
			name: "skips logging when event level is nil",
			event: JobEvent{
				JobSpec: "@daily",
				JobName: "rotate",
				Stage:   StageFinish,
			},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				capture = &slogCaptureHandler{}
				log     = slog.New(capture)
			)

			handler := NewSlogHandler(log).
				WithError(tc.levelers.error).
				WithEvent(tc.levelers.event)

			handler.Handle(tc.event)

			capture.assert(t, tc.expected)
		})
	}
}
