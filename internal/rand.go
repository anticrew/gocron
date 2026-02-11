package internal

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	RandNameDefaultSize = 8
)

func RandName(size int) string {
	if size < 1 {
		size = RandNameDefaultSize
	}

	b := make([]byte, size)
	// Не проверяем ошибку: в Go 1.25.0 rand.Read паникует при любой ошибке.
	_, _ = rand.Read(b)

	return base64.RawURLEncoding.EncodeToString(b)[:size]
}
