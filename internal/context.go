package internal

import (
	"context"
	"time"
)

type ContextFactory func(ctx context.Context) (context.Context, context.CancelFunc)

func TimeoutContextFactory(timeout time.Duration) ContextFactory {
	return func(ctx context.Context) (context.Context, context.CancelFunc) {
		return context.WithTimeout(ctx, timeout)
	}
}

func CancelContextFactory() ContextFactory {
	return context.WithCancel
}
