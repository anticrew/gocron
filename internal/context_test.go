package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeoutContextFactory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		timeout time.Duration
		spread  time.Duration
	}{
		{
			name:    "positive timeout",
			timeout: 150 * time.Millisecond,
			spread:  50 * time.Millisecond,
		},
		{
			name:    "negative timeout",
			timeout: -50 * time.Millisecond,
			spread:  50 * time.Millisecond,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			factory := TimeoutContextFactory(tc.timeout)
			require.NotNil(t, factory)

			ctx, cancel := factory(t.Context())
			require.NotNil(t, ctx)
			require.NotNil(t, cancel)
			t.Cleanup(cancel)

			deadline, ok := ctx.Deadline()
			assert.True(t, ok)

			var (
				wantMin = tc.timeout - tc.spread
				wantMax = tc.timeout + tc.spread
				actual  = time.Until(deadline)
			)

			assert.Greater(t, actual, wantMin)
			assert.Less(t, actual, wantMax)
		})
	}
}

func TestCancelContextFactory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "default",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			factory := CancelContextFactory()
			require.NotNil(t, factory)

			ctx, cancel := factory(t.Context())
			require.NotNil(t, ctx)
			require.NotNil(t, cancel)
			t.Cleanup(cancel)

			deadline, ok := ctx.Deadline()
			assert.False(t, ok)
			assert.Zero(t, deadline)
		})
	}
}
