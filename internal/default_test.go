package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDefaultNil(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		const (
			zeroValue    = 0
			defaultValue = 1
			validValue   = 15
		)

		defFn := func() int { return defaultValue }

		assert.Equal(t, defaultValue, WithDefault(zeroValue, defFn))
		assert.Equal(t, validValue, WithDefault(validValue, defFn))
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		const (
			zeroValue    = ""
			defaultValue = "string"
			validValue   = "message"
		)

		defFn := func() string { return defaultValue }
		assert.Equal(t, defaultValue, WithDefault(zeroValue, defFn))
		assert.Equal(t, validValue, WithDefault(validValue, defFn))
	})

	t.Run("context", func(t *testing.T) {
		t.Parallel()

		var (
			zeroValue    context.Context
			defaultValue = context.Background()
			validValue   = t.Context()
		)

		defFn := func() context.Context { return defaultValue }
		assert.Equal(t, defaultValue, WithDefault(zeroValue, defFn))
		assert.Equal(t, validValue, WithDefault(validValue, defFn))
	})
}
