package handlers

import (
	"encoding/json"
	"fmt"
	"iv_project/dto"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.SuccessResult{StatusCode: statusCode, Message: message, Data: data})
}

func ErrorResponse(w http.ResponseWriter, statusCode int, messages map[string]string, lang string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	message, exists := messages[lang]
	if !exists {
		message = messages["en"]
	}
	json.NewEncoder(w).Encode(dto.ErrorResult{StatusCode: statusCode, Message: message})
}

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

func CalculateDiscountedPrice(price uint, discountPercentage uint) uint {
	discount := (float64(discountPercentage) / 100) * float64(price)
	finalPrice := math.Round(float64(price) - discount)
	return uint(finalPrice)
}
