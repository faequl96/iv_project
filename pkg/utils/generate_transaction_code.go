package utils

import (
	"fmt"
	"time"
)

func GenerateTransactionCode() string {
	return fmt.Sprintf("TX-%s-%s", time.Now().Format("0601021504"), randomString(3))
}
