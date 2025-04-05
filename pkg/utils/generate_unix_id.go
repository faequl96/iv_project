package utils

import (
	"fmt"
	"time"
)

func GenerateUnixID() string {
	prefix := "USER"

	return fmt.Sprintf("%s-%s%s", prefix, time.Now().Format("060102"), randomString(5))
}
