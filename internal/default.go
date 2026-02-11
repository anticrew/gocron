package internal

func WithDefault[T comparable](val T, def func() T) T {
	var zero T
	if val == zero {
		return def()
	}

	return val
}
