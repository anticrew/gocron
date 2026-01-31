package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		size         int
		expectedSize int
	}{
		{
			name:         "zero size",
			size:         0,
			expectedSize: 8,
		},
		{
			name:         "non-zero size",
			size:         8,
			expectedSize: 8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			const attempts = 1000

			prev := make(map[string]struct{}, attempts)

			for range attempts {
				name := RandName(tc.size)
				assert.Len(t, name, tc.expectedSize)

				assert.NotContains(t, prev, name)
				prev[name] = struct{}{}
			}
		})
	}
}
