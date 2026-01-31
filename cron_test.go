package gocron

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCron_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		spec             string
		cmd              Cmd
		expectedErr      error
		expectedContains string
	}{
		{
			name:        "nil command returns ErrCommandIsNil",
			spec:        "* * * * * *",
			cmd:         nil,
			expectedErr: ErrCommandIsNil,
		},
		{
			name:             "invalid spec wraps cron error",
			spec:             "bad spec",
			cmd:              func(context.Context) error { return nil },
			expectedContains: "cron.AddJob",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewCron(t.Context()).Add(tc.spec, tc.cmd)

			switch {
			case tc.expectedErr != nil:
				assert.ErrorIs(t, err, tc.expectedErr)

			case len(tc.expectedContains) > 0:
				assert.ErrorContains(t, err, tc.expectedContains)
			}
		})
	}
}

func TestCron_DefaultHandler(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	var called bool
	c := NewCron(ctx, WithDefaultHandler(HandlerFunc(func(event JobEvent) {
		called = true
	})))

	c.MustAdd("@every 1s", func(context.Context) error { return nil })

	c.Start()

	time.Sleep(time.Second)

	require.NoError(t, c.Shutdown(ctx))

	assert.True(t, called)
}

func TestCron_DefaultTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		spec     string
		timeout  time.Duration
		job      time.Duration
		wait     time.Duration
		expected error
	}{
		{
			name:     "positive timeout sets a deadline",
			spec:     "@every 1s",
			timeout:  10 * time.Millisecond,
			job:      100 * time.Millisecond,
			wait:     time.Second,
			expected: context.DeadlineExceeded,
		},
		{
			name:    "non-positive timeout does not set a deadline",
			spec:    "@every 1s",
			timeout: 0,
			job:     100 * time.Millisecond,
			wait:    time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			c := NewCron(ctx, WithTimeout(tc.timeout))

			var err error
			c.MustAdd("* * * * * *", func(ctx context.Context) error {
				time.Sleep(tc.job)
				err = ctx.Err()
				return nil
			})

			c.Start()

			time.Sleep(tc.wait)

			require.NoError(t, c.Shutdown(ctx))

			assert.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestCron_Start(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	c := NewCron(ctx)

	var called int
	c.MustAdd("@every 1s", func(context.Context) error {
		called++
		return nil
	})

	for range 10 {
		c.Start()
		c.Start()
	}

	time.Sleep(time.Second)
	require.NoError(t, c.Shutdown(ctx))

	assert.Equal(t, 1, called)
}

func TestCron_Shutdown(t *testing.T) {
	t.Parallel()

	t.Run("before start", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c := NewCron(ctx)
		assert.ErrorIs(t, c.Shutdown(ctx), ErrCronNotRunning)
	})

	t.Run("several calls", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c := NewCron(ctx)
		c.Start()
		require.NoError(t, c.Shutdown(ctx))
		require.Error(t, c.Shutdown(ctx), ErrCronNotRunning)
	})
}

func TestCron_GracefulShutdown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		spec     string
		job      time.Duration
		wait     time.Duration
		timeout  time.Duration
		expected error
	}{
		{
			name:    "no timeout",
			spec:    "@every 1s",
			job:     3 * time.Second,
			wait:    10 * time.Second,
			timeout: time.Minute,
		},
		{
			name:     "timeout exceeded",
			spec:     "@every 1s",
			job:      3 * time.Second,
			wait:     10 * time.Second,
			timeout:  1 * time.Second,
			expected: context.DeadlineExceeded,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			var (
				started  atomic.Int32
				finished atomic.Int32
			)

			c := NewCron(ctx)
			c.MustAdd(tc.spec, func(ctx context.Context) error {
				started.Add(1)

				time.Sleep(tc.job)

				if err := ctx.Err(); err != nil {
					return err
				}

				finished.Add(1)
				return nil
			})

			c.Start()

			time.Sleep(tc.wait)

			shutdownCtx, cancel := context.WithTimeout(ctx, tc.timeout)
			t.Cleanup(cancel)

			require.ErrorIs(t, c.Shutdown(shutdownCtx), tc.expected)

			assert.NotZero(t, started.Load())
			if tc.expected == nil {
				assert.Equal(t, started.Load(), finished.Load())
			}
		})
	}
}

func TestCronMustAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec string
		cmd  Cmd
	}{
		// Verifies MustAdd panics on invalid input to match the documented contract.
		{
			name: "must add panics on invalid spec",
			spec: "bad spec",
			cmd:  func(context.Context) error { return nil },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cron := NewCron(t.Context())

			assert.Panics(t, func() {
				cron.MustAdd(tc.spec, tc.cmd)
			})
		})
	}
}
