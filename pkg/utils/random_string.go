package utils

import (
	"math/rand"
	"time"
)

func randomString(n int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
