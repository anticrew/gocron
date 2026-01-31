package internal

import (
	"context"
	"sync"
)

func Wait(ctx context.Context, wg *sync.WaitGroup) error {
	ch := make(chan struct{}, 1)

	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}
