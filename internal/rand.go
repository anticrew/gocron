package internal

import (
	"crypto/rand"
	"encoding/base64"
)

func RandName(size int) string {
	b := make([]byte, size)
	_, _ = rand.Read(b)

	return base64.RawURLEncoding.EncodeToString(b)
}
