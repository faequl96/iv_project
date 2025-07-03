package utils

import (
	"fmt"
	"time"
)

func GenerateReferenceNumber(paymentMethod string) string {
	prefix := map[string]string{
		"iv_coin":       "IVC",
		"gopay":         "GOP",
		"qris":          "QRI",
		"bank_transfer": "BTF",
	}[paymentMethod]
	if prefix == "" {
		prefix = "TXN"
	}

	return fmt.Sprintf("%s-%s%s", prefix, time.Now().Format("060102"), randomString(5))
}
