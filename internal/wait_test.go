package internal

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWait(t *testing.T) {
	t.Parallel()

	const (
		maxTestDuration = 15 * time.Second
		timeoutDuration = 1 * time.Second
		sleepDuration   = 5 * time.Second
	)

	require.Greater(t, sleepDuration, timeoutDuration)
	require.Greater(t, maxTestDuration, sleepDuration)

	tests := []struct {
		name           string
		timeout        time.Duration
		run            func(wg *sync.WaitGroup)
		expectedFinish bool
		expectedError  error
	}{
		{
			name:    "zero deadline exceeded",
			timeout: 0,
			run: func(wg *sync.WaitGroup) {
				defer wg.Done()
			},
			expectedError: context.DeadlineExceeded,
		},
		{
			name:    "timeout deadline exceeded",
			timeout: timeoutDuration,
			run: func(wg *sync.WaitGroup) {
				defer wg.Done()
				time.Sleep(sleepDuration)
			},
			expectedError: context.DeadlineExceeded,
		},
		{
			name:    "wait",
			timeout: timeoutDuration,
			run: func(wg *sync.WaitGroup) {
				defer wg.Done()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testCtx, cancel := context.WithTimeout(t.Context(), maxTestDuration)
			t.Cleanup(cancel)

			ctx, cancel := context.WithTimeout(testCtx, tc.timeout)
			t.Cleanup(cancel)

			var (
				wg      = &sync.WaitGroup{}
				errChan = make(chan error, 1)
			)

			wg.Add(1)
			go tc.run(wg)
			func() {
				errChan <- Wait(ctx, wg)
				close(errChan)
			}()

			select {
			case err := <-errChan:
				require.ErrorIs(t, err, tc.expectedError)

			case <-testCtx.Done():
				t.Fatal("timed out")
			}
		})
	}
}
