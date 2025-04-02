package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateReferenceNumber(paymentMethod string) string {
	prefix := map[string]string{
		"iv_coin":         "IVC",
		"manual_transfer": "MTF",
		"auto_transfer":   "ATF",
		"gopay":           "GOP",
	}[paymentMethod]
	if prefix == "" {
		prefix = "TXN"
	}

	return fmt.Sprintf("%s-%s%s", prefix, time.Now().Format("060102"), randomString(5))
}

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
