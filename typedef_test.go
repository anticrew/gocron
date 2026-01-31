package gocron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerFuncHandle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		event JobEvent
	}{
		{
			name: "pass without changes",
			event: JobEvent{
				JobSpec: "@every 1s",
				JobName: "cleanup",
				Stage:   StageExec,
				Error:   assert.AnError,
			},
		},
		{
			name: "pass no error",
			event: JobEvent{
				JobSpec: "0 0 * * *",
				JobName: "daily",
				Stage:   StageFinish,
				Error:   nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got JobEvent
			handler := HandlerFunc(func(event JobEvent) {
				got = event
			})

			handler.Handle(tc.event)

			assert.Equal(t, tc.event, got)
		})
	}
}
