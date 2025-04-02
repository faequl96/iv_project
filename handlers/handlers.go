package handlers

import (
	"encoding/json"
	"iv_project/dto"
	"net/http"
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
