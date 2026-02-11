package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  int
		err  error
	}{
		{
			name: "no error",
			val:  42,
		},
		{
			name: "panic on error",
			val:  7,
			err:  assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fn := func() {
				val := Must(tc.val, tc.err)
				assert.Equal(t, tc.val, val)
			}

			if tc.err == nil {
				assert.NotPanics(t, fn)
			} else {
				assert.PanicsWithError(t, tc.err.Error(), fn)
			}
		})
	}
}
