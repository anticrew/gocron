package gocron

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type jobLock struct {
	lockErr   error
	unlockErr error
}

func (l *jobLock) Lock(context.Context) error {
	return l.lockErr
}

func (l *jobLock) Unlock(context.Context) error {
	return l.unlockErr
}

type jobHandler struct {
	events []JobEvent
}

func (h *jobHandler) Handle(event JobEvent) {
	h.events = append(h.events, event)
}

func TestJob_Timeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		timeout       time.Duration
		job           time.Duration
		expectedError error
	}{
		{
			name:          "positive timeout creates deadline",
			timeout:       5 * time.Millisecond,
			job:           time.Second,
			expectedError: context.DeadlineExceeded,
		},
		{
			name:    "non-positive timeout disables deadline",
			timeout: 0,
			job:     time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var called bool
			j := newJob(t.Context(), "spec", func(ctx context.Context) error {
				called = true
				time.Sleep(tc.job)
				return ctx.Err()
			})
			j.WithTimeout(tc.timeout)
			j.WithHandler(HandlerFunc(func(event JobEvent) {
				if event.Stage != StageExec {
					return
				}

				assert.Equal(t, tc.expectedError, event.Error)
			}))

			j.Run()

			assert.True(t, called)
		})
	}
}

func TestJob_Name(t *testing.T) {
	t.Parallel()

	const (
		spec = "spec"
		name = "name"
	)

	j := newJob(t.Context(), spec, func(context.Context) error {
		return nil
	})
	j.WithName(name)

	var actualName string
	j.WithHandler(HandlerFunc(func(event JobEvent) {
		actualName = event.JobName
	}))

	j.Run()

	assert.Equal(t, name, actualName)
}

func TestJob_Handler(t *testing.T) {
	t.Parallel()

	const (
		spec = "spec"
		name = "name"
	)

	tests := []struct {
		name     string
		err      error
		lock     Lock
		expected []JobEvent
	}{
		{
			name: "start only",
			err:  nil,
			lock: &jobLock{
				lockErr: assert.AnError,
			},
			expected: []JobEvent{
				{
					JobSpec: spec,
					JobName: name,
					Stage:   StageStart,
					Error:   assert.AnError,
				},
			},
		},
		{
			name: "start execute",
			err:  assert.AnError,
			expected: []JobEvent{
				{
					JobSpec: spec,
					JobName: name,
					Stage:   StageStart,
					Error:   nil,
				},
				{
					JobSpec: spec,
					JobName: name,
					Stage:   StageExec,
					Error:   assert.AnError,
				},
				{
					JobSpec: spec,
					JobName: name,
					Stage:   StageFinish,
					Error:   nil,
				},
			},
		},
		{
			name: "start execute finish",
			err:  nil,
			lock: &jobLock{
				unlockErr: assert.AnError,
			},
			expected: []JobEvent{
				{
					JobSpec: spec,
					JobName: name,
					Stage:   StageStart,
					Error:   nil,
				},
				{
					JobSpec: spec,
					JobName: name,
					Stage:   StageExec,
					Error:   nil,
				},
				{
					JobSpec: spec,
					JobName: name,
					Stage:   StageFinish,
					Error:   assert.AnError,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			j := newJob(t.Context(), spec, func(_ context.Context) error {
				return tc.err
			})

			events := make([]JobEvent, 0, 3)
			j.WithHandler(HandlerFunc(func(event JobEvent) {
				events = append(events, event)
			})).
				WithLock(tc.lock).
				WithName(name)

			j.Run()

			assert.Equal(t, tc.expected, events)
		})
	}

	t.Run("no handler", func(t *testing.T) {
		t.Parallel()

		var called bool
		j := newJob(t.Context(), "spec", func(context.Context) error {
			called = true
			return assert.AnError
		})

		assert.NotPanics(t, j.Run)
		assert.True(t, called)
	})
}

func TestJob_WaitGroup(t *testing.T) {
	t.Parallel()

	const (
		sleep  = time.Second
		spread = 50 * time.Millisecond
	)

	j := newJob(t.Context(), "spec", func(context.Context) error {
		time.Sleep(sleep)
		return nil
	})

	wg := &sync.WaitGroup{}
	j.withWaitGroup(wg)

	start := time.Now()
	go j.Run()

	wg.Wait()

	assert.GreaterOrEqual(t, time.Since(start), sleep-spread)
}
